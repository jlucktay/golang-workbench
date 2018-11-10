package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":8080", "Listen address")
	dir    = flag.String("dir", ".", "Directory to serve")
)

func main() {
	flag.Parse()
	log.Printf("Listening on '%q'...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir))))
}
