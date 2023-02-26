// Chaining things together through the magic of interfaces!
package main

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	const input = `H4sIAAAAAAAA/3qyd8GT3WueLt37fk/Ps46JT/d1vFw942nrlqezFzzZsQskMmfFi/0zX7b3cAECAAD//0G6Zq8rAAAA`
	var r io.Reader = strings.NewReader(input)
	r = base64.NewDecoder(base64.StdEncoding, r)
	r, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(os.Stdout, r)
}
