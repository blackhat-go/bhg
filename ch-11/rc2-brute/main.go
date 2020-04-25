package main

import (
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"log"
	"regexp"
	"sync"

	luhn "github.com/joeljunstrom/go-luhn"

	"github.com/blackhat-go/bhg/ch-11/rc2-brute/rc2"
)

var numeric = regexp.MustCompile(`^\d{8}$`)

type CryptoData struct {
	block cipher.Block
	key   []byte
}

func generate(start, stop uint64, out chan<- *CryptoData, done <-chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		var (
			block cipher.Block
			err   error
			key   []byte
			data  *CryptoData
		)
		for i := start; i <= stop; i++ {
			key = make([]byte, 8)
			select {
			case <-done:
				return
			default:
				binary.BigEndian.PutUint64(key, i)
				if block, err = rc2.New(key[3:], 40); err != nil {
					log.Fatalln(err)
				}
				data = &CryptoData{
					block: block,
					key:   key[3:],
				}
				out <- data
			}
		}
	}()

	return
}

func decrypt(ciphertext []byte, in <-chan *CryptoData, done chan struct{}, wg *sync.WaitGroup) {
	size := rc2.BlockSize
	plaintext := make([]byte, len(ciphertext))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range in {
			select {
			case <-done:
				return
			default:
				data.block.Decrypt(plaintext[:size], ciphertext[:size])
				if numeric.Match(plaintext[:size]) {
					data.block.Decrypt(plaintext[size:], ciphertext[size:])
					if luhn.Valid(string(plaintext)) && numeric.Match(plaintext[size:]) {
						log.Printf("Card [%s] found using key [%x]\n", plaintext, data.key)
						close(done)
						return
					}
				}
			}
		}
	}()
}

func main() {
	var (
		err        error
		ciphertext []byte
	)
	/*
	   Assume ECB mode
	   Assume no padding needed
	   key        = e612d0bbb6
	   ciphertext = 0986f2cc1ebdc5c2e25d04a136fa1a6b
	   plaintext  = 4532651325506680
	*/
	if ciphertext, err = hex.DecodeString("0986f2cc1ebdc5c2e25d04a136fa1a6b"); err != nil {
		log.Fatalln(err)
	}

	var prodWg, consWg sync.WaitGroup
	var min, max, prods = uint64(0x0000000000), uint64(0xffffffffff), uint64(75)
	var step = (max - min) / prods

	done := make(chan struct{})
	work := make(chan *CryptoData, 100)
	if (step * prods) < max {
		step += prods
	}
	var start, end = min, min + step
	log.Println("Starting producers...")
	for i := uint64(0); i < prods; i++ {
		if end > max {
			end = max
		}
		generate(start, end, work, done, &prodWg)
		end += step
		start += step
	}
	log.Println("Producers started!")
	log.Println("Starting consumers...")
	for i := 0; i < 30; i++ {
		decrypt(ciphertext, work, done, &consWg)
	}
	log.Println("Consumers started!")
	log.Println("Now we wait...")
	prodWg.Wait()
	close(work)
	consWg.Wait()
	log.Println("Brute-force complete")
}
