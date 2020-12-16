package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	client := http.Client{
		Timeout: time.Minute * 3,
	}

	// change the address to match the common name of the certificate
	resp, err := client.Get("https://example.test:9090")
	if err != nil {
		log.Fatalf("error making get request: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}
	fmt.Println(string(body))
}
