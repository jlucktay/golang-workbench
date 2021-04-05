package main

import (
	"bytes"
	"flag"
	"fmt"
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
	mux.HandleFunc("/tokensignin", tokenSignIn)

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

func tokenSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp := fmt.Sprintf("Method Not Allowed: %s", r.Method)
		http.Error(w, resp, http.StatusMethodNotAllowed)
		log.Println(resp)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("could not parse request form: %v", err)
		return
	}

	idToken, tokenPresent := r.Form["idtoken"]
	if !tokenPresent {
		resp := "Bad Request: no 'idtoken' in form"
		http.Error(w, resp, http.StatusBadRequest)
		log.Println(resp)
		return
	}

	if len(idToken) != 1 {
		resp := "Bad Request: idtoken slice contains incorrect numnber of elements"
		http.Error(w, resp, http.StatusBadRequest)
		log.Println(resp)
		return
	}

	log.Printf("ID token: '%s'", idToken[0])

	fmt.Println("TODO: Verify the integrity of the ID token " +
		"(https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token)")
}
