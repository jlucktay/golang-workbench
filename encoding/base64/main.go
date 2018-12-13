package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	serialEncoded := "ASNHQPdfYCLd"
	serialDecoded, errDecode := base64.StdEncoding.DecodeString(serialEncoded)
	if errDecode != nil {
		log.Fatal(errDecode)
	}
	fmt.Println("serialDecoded:", serialDecoded, string(serialDecoded))
	fmt.Printf("%%s: %s\n", serialDecoded)
}
