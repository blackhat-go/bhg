package metadata

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var refRegex = regexp.MustCompile(`[0-9]+ [0-9]+ R`)

type PDFBytes []byte

type Reference struct {
	ObjectID int
	GenID    int
}

type Trailer struct {
	Root *Reference
	Info *Reference
	Prev int
}

type Info struct {
	XMLName  xml.Name `xml:"xmpmeta"`
	Author   string   `xml:"RDF>Description>creator"`
	Creator  string   `xml:"RDF>Description>CreatorTool"`
	Producer string   `xml:"RDF>Description>Producer"`
}

type Meta struct {
}

type XRefObject struct {
	ObjectID int
	Offset   int64
}

type XRef struct {
	StartID   int
	Count     int
	ObjectRef []XRefObject
}

func NewPDFData(buf []byte, stripNewLines bool) PDFBytes {
	var ret PDFBytes
	b := bytes.Trim(buf, "\x20\x09\x00\x0C")
	if stripNewLines {
		b = bytes.Replace(b, []byte("\x0A"), []byte{}, -1)
		b = bytes.Replace(b, []byte("\x0D"), []byte{}, -1)
	}
	ret = PDFBytes(b)
	return ret
}

func NewPropertiesFromPDFDoc(file string) (info []Info, err error) {
	var xref *XRef
	var trailer *Trailer
	var buf []byte
	info = make([]Info, 0)

	buf, err = ioutil.ReadFile(file)
	if err != nil {
		return
	}

	current := -1
	next := -1
	var iRef *Reference
	var rRef *Reference

	for {
		next, xref, trailer, err = ParseFileTrailer(buf, current)
		if err != nil {
			return
		}
		if xref == nil {
			if next < 0 {
				err = errors.New("No next xref section")
				return
			}
			current = next
			continue
		}
		if trailer == nil {
			err = errors.New("Missing Trailer")
			return
		}

		if rRef == nil {
			rRef = trailer.Root

			d := xref.FetchData(buf, *rRef)
			if len(d) == 0 {
				err = errors.New("Cannot fetch Root data")
				return
			}

			if d.TypeOf() != "MAP" {
				return nil, errors.New("Unexpected Root Type")
			}

			iRef, err = d.GetMetaRef()
			if err != nil {
				return
			}
		}

		if iRef != nil {
			if iRef.ObjectID < xref.StartID || iRef.ObjectID > xref.StartID+xref.Count {
				current = next
				continue
			}
			d := xref.FetchData(buf, *iRef)
			s := d.ToXMLStream()
			var i Info
			if err = xml.Unmarshal(s, &i); err != nil {
				return
			}
			info = append(info, i)
		}

		break
	}

	current = -1
	next = -1
	for {
		next, xref, trailer, err = ParseFileTrailer(buf, current)
		if err != nil {
			return
		}
		if xref == nil {
			if next < 0 {
				err = errors.New("No next xref section")
				return
			}
			current = next
			continue
		}
		if trailer == nil {
			err = errors.New("Missing Trailer")
			return
		}
		if trailer.Info != nil {
			iRef = trailer.Info
		}

		d := xref.FetchData(buf, *iRef)
		if len(d) == 0 {
			current = next
			continue
		}
		if d.TypeOf() != "MAP" {
			return nil, errors.New("Unexpected Trailer Info Type")
		}
		i, err := d.ToInfo(*xref, buf)
		if err != nil {
			return nil, err
		}
		info = append(info, *i)

		break
	}

	return info, nil
}

func ParseFileTrailer(buf []byte, current int) (nextOffset int, xref *XRef, trailer *Trailer, err error) {
	nextOffset = -1

	xrefOffset := current
	if xrefOffset < 0 {
		startXrefOffset := bytes.LastIndex(buf, []byte("startxref"))
		eofOffset := bytes.LastIndex(buf, []byte("%%EOF"))
		b := buf[startXrefOffset+len("startxref") : eofOffset]
		d := NewPDFData(b, true)
		xrefOffset, err = strconv.Atoi(string(d))
		if err != nil {
			return
		}
	}

	trailerOffset := bytes.Index(buf[xrefOffset:], []byte("trailer")) + xrefOffset
	d := NewPDFData(buf[xrefOffset+len("xref"):trailerOffset], false)
	xref, err = d.ToXRef()
	if err != nil {
		return
	}

	startxref := bytes.Index(buf[xrefOffset:], []byte("startxref")) + xrefOffset
	d = NewPDFData(buf[trailerOffset+len("trailer"):startxref], true)
	trailer, err = d.ToTrailer()
	if err != nil {
		return
	}
	nextOffset = trailer.Prev
	return

}

func (b PDFBytes) TypeOf() string {
	buf := bytes.Trim(b, "\x0A\x0D\x20\x09\x00\x0C")
	if bytes.Index(buf, []byte("<<")) == 0 && bytes.LastIndex(buf, []byte(">>")) == len(buf)-2 {
		return "MAP"
	}
	if bytes.Index(buf, []byte("(")) == 0 && bytes.LastIndex(buf, []byte(")")) == len(buf)-1 {
		return "STRING"
	}
	if refRegex.Match(buf) {
		return "REF"
	}
	if _, err := strconv.Atoi(string(buf)); err == nil {
		return "INT"
	}
	return ""
}

func (b PDFBytes) ToTrailer() (*Trailer, error) {
	var trailer Trailer
	buf := []byte(b)
	if b.TypeOf() != "MAP" {
		return nil, errors.New("Invalid Trailer Record Format")
	}
	buf = bytes.Replace(buf, []byte("<"), []byte{}, -1)
	buf = bytes.Replace(buf, []byte(">"), []byte{}, -1)
	buf = bytes.Trim(buf, "\x0A\x0D\x20\x09\x00\x0C")
	fields := bytes.Split(buf, []byte("/"))
	for _, field := range fields {
		tokens := bytes.Split(field, []byte(" "))
		if bytes.Index(field, []byte("Root")) == 0 {
			oid, err := strconv.Atoi(string(tokens[1]))
			if err != nil {
				return nil, err
			}
			gid, err := strconv.Atoi(string(tokens[2]))
			if err != nil {
				return nil, err
			}
			trailer.Root = &Reference{ObjectID: oid, GenID: gid}
		}
		if bytes.Index(field, []byte("Info")) == 0 {
			oid, err := strconv.Atoi(string(tokens[1]))
			if err != nil {
				return nil, err
			}
			gid, err := strconv.Atoi(string(tokens[2]))
			if err != nil {
				return nil, err
			}
			trailer.Info = &Reference{ObjectID: oid, GenID: gid}
		}
		if bytes.Index(field, []byte("Prev")) == 0 {
			t := strings.Trim(string(field[len("Prev"):]), "\x0A\x0D\x20\x09\x00\x0C")
			p, err := strconv.Atoi(t)
			if err != nil {
				return nil, err
			}
			trailer.Prev = p
		}
	}
	return &trailer, nil
}

func (b PDFBytes) ToInfo(xref XRef, doc []byte) (*Info, error) {
	var info Info
	buf := []byte(b)
	if b.TypeOf() != "MAP" {
		return nil, errors.New("Invalid Info Record Format")
	}
	buf = bytes.Replace(buf, []byte("<"), []byte{}, -1)
	buf = bytes.Replace(buf, []byte(">"), []byte{}, -1)
	buf = bytes.Trim(buf, "\x0A\x0D\x20\x09\x00\x0C")
	fields := bytes.Split(buf, []byte("/"))
	for _, field := range fields {
		if bytes.Index(field, []byte("Author")) == 0 {
			data := NewPDFData(field[len("Author"):], true)
			if data.TypeOf() == "STRING" {
				info.Author = data.ToString()
				continue
			}
			if data.TypeOf() == "REF" {
				tokens := bytes.Split(data, []byte(" "))
				oid, err := strconv.Atoi(string(tokens[0]))
				if err != nil {
					return nil, err
				}
				gid, err := strconv.Atoi(string(tokens[1]))
				if err != nil {
					return nil, err
				}
				ref := Reference{ObjectID: oid, GenID: gid}
				d := xref.FetchData(doc, ref)
				info.Author = d.ToString()
			}
		}
		if bytes.Index(field, []byte("Creator")) == 0 {
			data := NewPDFData(field[len("Creator"):], true)
			if data.TypeOf() == "STRING" {
				info.Creator = data.ToString()
				continue
			}
			if data.TypeOf() == "REF" {
				tokens := bytes.Split(data, []byte(" "))
				oid, err := strconv.Atoi(string(tokens[0]))
				if err != nil {
					return nil, err
				}
				gid, err := strconv.Atoi(string(tokens[1]))
				if err != nil {
					return nil, err
				}
				ref := Reference{ObjectID: oid, GenID: gid}
				d := xref.FetchData(doc, ref)
				info.Creator = d.ToString()
			}
		}
		if bytes.Index(field, []byte("Producer")) == 0 {
			data := NewPDFData(field[len("Producer"):], true)
			if data.TypeOf() == "STRING" {
				info.Producer = data.ToString()
				continue
			}
			if data.TypeOf() == "REF" {
				tokens := bytes.Split(data, []byte(" "))
				oid, err := strconv.Atoi(string(tokens[0]))
				if err != nil {
					return nil, err
				}
				gid, err := strconv.Atoi(string(tokens[1]))
				if err != nil {
					return nil, err
				}
				ref := Reference{ObjectID: oid, GenID: gid}
				d := xref.FetchData(doc, ref)
				info.Producer = d.ToString()
			}
		}
	}
	return &info, nil
}

func (b PDFBytes) GetMetaRef() (ref *Reference, err error) {
	var oid, gid int
	buf := []byte(b)
	if b.TypeOf() != "MAP" {
		err = errors.New("Invalid Root Record Format")
		return
	}
	buf = bytes.Replace(buf, []byte("<"), []byte{}, -1)
	buf = bytes.Replace(buf, []byte(">"), []byte{}, -1)
	buf = bytes.Trim(buf, "\x0A\x0D\x20\x09\x00\x0C")
	fields := bytes.Split(buf, []byte("/"))
	for _, field := range fields {
		if bytes.Index(field, []byte("Metadata")) == 0 {
			data := NewPDFData(field[len("Metadata"):], true)
			tokens := bytes.Split(data, []byte(" "))
			oid, err = strconv.Atoi(string(tokens[0]))
			if err != nil {
				return
			}
			gid, err = strconv.Atoi(string(tokens[1]))
			if err != nil {
				return
			}
			ref = &Reference{ObjectID: oid, GenID: gid}
			break
		}
	}
	return
}

func (b PDFBytes) ToXRef() (*XRef, error) {
	var xref XRef
	recLength := 20
	if len(b) < recLength {
		return nil, nil
	}
	headerLength := len(b) % recLength
	header := bytes.Trim(b[:headerLength], "\x0A\x0D\x20\x09\x00\x0C")
	tokens := bytes.Split(header, []byte(" "))
	id, err := strconv.Atoi(string(tokens[0]))
	if err != nil {
		return nil, err
	}
	xref.StartID = id

	count, err := strconv.Atoi(string(tokens[1]))
	if err != nil {
		return nil, err
	}
	xref.Count = count

	offset := headerLength
	xref.ObjectRef = make([]XRefObject, 0, xref.Count)
	for i := xref.StartID; i < xref.StartID+xref.Count; i++ {
		x := b[offset : offset+recLength]
		tokens := bytes.Split(x, []byte(" "))
		objectOffset, err := strconv.Atoi(string(tokens[0]))
		if err != nil {
			return nil, err
		}
		xrefObj := &XRefObject{Offset: int64(objectOffset), ObjectID: i}
		xref.ObjectRef = append(xref.ObjectRef, *xrefObj)
		offset += 20
	}

	return &xref, nil
}

func (b PDFBytes) ToXMLStream() []byte {
	return b[bytes.Index(b, []byte("stream"))+len("stream") : bytes.LastIndex(b, []byte("endstream"))]
}

func (b PDFBytes) ToString() string {
	return string(bytes.Trim(b, "\x0A\x0D\x20\x09\x00\x0C()"))
}

func (x XRef) FetchData(buf []byte, ref Reference) PDFBytes {
	var ret PDFBytes
	for _, obj := range x.ObjectRef {
		if obj.ObjectID == ref.ObjectID {
			b := buf[obj.Offset:]
			b = b[bytes.Index(b, []byte("obj"))+len("obj") : bytes.Index(b, []byte("endobj"))]
			ret = NewPDFData(b, true)
			break
		}
	}
	return ret
}
