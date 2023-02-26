package readbyte

import (
	"fmt"
	"io"
	"log"
)

func readByte() {
	err := io.EOF // force an error
	if err != nil {
		fmt.Println("ERROR")
		log.Print("Couldn't read first byte")

		return
	}
}
