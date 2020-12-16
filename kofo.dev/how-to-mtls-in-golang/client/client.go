package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	cert, err := ioutil.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	certificate, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
				RootCAs:      caCertPool,
			},
		},
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
