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

	// Lock 'em in
	flag.Parse()

	// Prepare the login page template
	tpl := template.Must(template.New("gsifw.html").ParseFiles("gsifw.html"))

	if *clientID == "" {
		log.Fatal("missing Google Client ID; set GOOGLE_CLIENT_ID in env or '--client-id' flag")
	}

	rootPage, err := prepareGSIFWBytes(tpl, *clientID)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("/", googleSignInForWebsites(rootPage))

	log.Printf("server listening on '%s'...", *address)
	log.Fatal(http.ListenAndServe(*address, mux))
}

// prepareGSIFWBytes will execute the given template to render the clientID into place, and return a byte-slice
// representation of the root page.
func prepareGSIFWBytes(tpl *template.Template, clientID string) ([]byte, error) {
	data := struct{ ClientID string }{ClientID: clientID}

	b := &bytes.Buffer{}
	if err := tpl.Execute(b, data); err != nil {
		log.Printf("could not execute template into buffer: %v", err)
		return nil, err
	}

	return gohtml.FormatBytes(b.Bytes()), nil
}

// googleSignInForWebsites runs a static-ish page through html/template and serves it.
func googleSignInForWebsites(page []byte) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)

		if _, err := w.Write(page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("could not write page bytes to ResponseWriter: %v", err)
			return
		}
	}
}
