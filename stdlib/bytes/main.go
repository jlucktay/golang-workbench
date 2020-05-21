package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var b bytes.Buffer // A Buffer needs no initialization.
	b.Write([]byte("Hello "))
	fmt.Fprintf(&b, "world!\n")
	_, _ = b.WriteTo(os.Stdout)

	fmt.Printf("'%s'\n", b.String())
	fmt.Printf("'%s'\n", b.String())

	// Second example with an io.TeeReader
	var c bytes.Buffer
	c.WriteString("Another buffer...\n")
	t := io.TeeReader(&c, os.Stdout)

	fmt.Println("before ReadAll(t)")
	xc, _ := ioutil.ReadAll(t)
	fmt.Println("after ReadAll(t)")

	fmt.Printf("Contents of 'xc' as string: '%s'\n", xc)
	fmt.Printf("Re-read buffer 'c': '%s'\n", c.String())

	// Third example which is really convoluted and also features a strings.Builder
	var d bytes.Buffer
	d.WriteRune('h')
	d.WriteRune('e')
	d.WriteRune('l')
	d.WriteRune('l')
	d.WriteRune('o')

	sb := strings.Builder{}

	fmt.Printf("sb:          '%#v'\n", sb)
	fmt.Printf("sb.String(): '%#v'\n", sb.String())

	t2 := io.TeeReader(&d, &sb)
	t2p := make([]byte, 1)

	for {
		n, errRead := t2.Read(t2p)
		if errRead != nil {
			if errRead == io.EOF {
				break
			}

			log.Fatal(errRead)
		}

		fmt.Printf("n:           '%d'\n", n)
		fmt.Printf("sb:          '%#v'\n", sb)
		fmt.Printf("sb.String(): '%s'\n", sb.String())
		fmt.Printf("sb.String(): '%s'\n", sb.String())
	}
}
