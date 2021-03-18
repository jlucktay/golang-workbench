package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/yosssi/gohtml"
)

func main() {
	// default credential flag to env var
	clientID := flag.String("client-id", os.Getenv("GOOGLE_CLIENT_ID"), "Google Client ID")

	// default address to localhost for development
	address := flag.String("server-address", "localhost:8080", "Server address to listen on")
	flag.Parse()

	if *clientID == "" {
		log.Fatal("missing Google Client ID; set GOOGLE_CLIENT_ID in env or '--client-id' flag")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", googleSignInForWebsites(*clientID))

	log.Printf("server listening on '%s'...\n", *address)
	log.Fatal(http.ListenAndServe(*address, mux))
}

// googleSignInForWebsites runs a static-ish page through html/template and serves it.
func googleSignInForWebsites(clientID string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		tpl := template.Must(template.New("gsifw.html").ParseFiles("gsifw.html"))
		data := struct{ ClientID string }{ClientID: clientID}

		b := &bytes.Buffer{}
		if err := tpl.Execute(b, data); err != nil {
			log.Printf("could not execute template into buffer: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fBytes := gohtml.Format(b.String())
		if _, err := w.Write([]byte(fBytes)); err != nil {
			log.Printf("could not write formatted bytes to ResponseWriter: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
