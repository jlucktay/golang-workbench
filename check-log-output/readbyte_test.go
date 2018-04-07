package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func readByte() {
	err := io.EOF // force an error

	if err != nil {
		fmt.Println("ERROR")
		log.Print("Couldn't read first byte")

		return
	}
}

func TestReadByte(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	defer func() {
		log.SetOutput(os.Stderr)
	}()

	readByte()
	t.Log(buf.String())
}
