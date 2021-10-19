package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Incrementer struct {
	mx      sync.Mutex
	counter uint64
}

func (i *Incrementer) String() string {
	i.mx.Lock()
	defer i.mx.Unlock()

	i.counter++

	return fmt.Sprintf("%d", i.counter)
}

func main() {
	var inc Incrementer

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "")
		http.Redirect(w, r, "/counter", http.StatusMovedPermanently)
	})

	http.HandleFunc("/counter", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, inc.String())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
