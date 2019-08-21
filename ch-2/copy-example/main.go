package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// FooReader defines an io.Reader to read from stdin.
type FooReader struct{}

// Read reads data from stdin.
func (fooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

// FooWriter defines an io.Writer to write to Stdout.
type FooWriter struct{}

// Write writes data to Stdout.
func (fooWriter *FooWriter) Write(b []byte) (int, error) {
	fmt.Print("out> ")
	return os.Stdout.Write(b)
}

func main() {
	// Instantiate reader and writer.
	var (
		reader FooReader
		writer FooWriter
	)

	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write data")
	}
}
