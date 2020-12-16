package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// set up handler to listen to root path
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("new %s request\n", r.Method)
		fmt.Fprintf(w, "hello world\n")
	})

	// load CA certificate file and add it to list of client CAs
	caCertFile, err := ioutil.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertFile)

	// serve on port 9090 of local host
	server := http.Server{
		Addr:    ":9090",
		Handler: handler,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  certPool,
			MinVersion: tls.VersionTLS12,
		},
	}

	// serve the endpoint with tls encryption
	if err := server.ListenAndServeTLS("../certs/server.crt", "../certs/server.key"); err != nil {
		log.Fatalf("error listening to port: %v", err)
	}
}
