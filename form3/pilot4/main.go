package main

import (
	"log"
	"net/http"
)

func main() {
	a := newApiServer()
	log.Fatal(http.ListenAndServe(":8080", a.router))
}
