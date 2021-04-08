package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yosssi/gohtml"
	"google.golang.org/api/idtoken"
)

const audienceSuffix = ".apps.googleusercontent.com"

var audience string

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

	audience = *clientID

	if strings.HasSuffix(*clientID, audienceSuffix) {
		*clientID = strings.TrimSuffix(*clientID, audienceSuffix)
	} else {
		audience += audienceSuffix
	}

	rootPage, err := prepareGSIFWBytes(tpl, *clientID)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/favicon.ico", chi.NewMux().NotFoundHandler())
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write(rootPage); err != nil {
			resp := fmt.Errorf("%s: could not write page bytes to ResponseWriter: %w",
				http.StatusText(http.StatusInternalServerError), err)
			http.Error(w, resp.Error(), http.StatusInternalServerError)
			log.Println(resp)
			return
		}
	})
	r.Post("/tokensignin", tokenSignIn)

	log.Printf("server listening on '%s'...", *address)
	log.Fatal(http.ListenAndServe(*address, r))
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

func tokenSignIn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		resp := fmt.Errorf("%s: could not parse request form: %w", http.StatusText(http.StatusBadRequest), err)
		http.Error(w, resp.Error(), http.StatusBadRequest)
		log.Println(resp)
		return
	}

	idToken, tokenPresent := r.Form["idtoken"]
	if !tokenPresent {
		resp := fmt.Sprintf("%s: no 'idtoken' in form", http.StatusText(http.StatusBadRequest))
		http.Error(w, resp, http.StatusBadRequest)
		log.Println(resp)
		return
	}

	if len(idToken) != 1 {
		resp := fmt.Sprintf("%s: idtoken slice contains incorrect number of elements",
			http.StatusText(http.StatusBadRequest))
		http.Error(w, resp, http.StatusBadRequest)
		log.Println(resp)
		return
	}

	idtp, err := verifyIntegrity(idToken[0])
	if err != nil {
		resp := fmt.Errorf("%s: could not verify integrity of the ID token: %w",
			http.StatusText(http.StatusBadRequest), err)
		http.Error(w, resp.Error(), http.StatusBadRequest)
		log.Println(resp)
		return
	}

	emailVerified, evOK := idtp.Claims["email_verified"]
	if !evOK {
		return
	}

	if bEmailVerified, ok := emailVerified.(bool); !ok || !bEmailVerified {
		return
	}

	email, ok := idtp.Claims["email"]
	if !ok {
		return
	}

	sEmail, ok := email.(string)
	if !ok || len(sEmail) == 0 {
		return
	}

	if _, err := w.Write([]byte(sEmail)); err != nil {
		resp := fmt.Errorf("%s: could not write bytes to ResponseWriter: %w",
			http.StatusText(http.StatusInternalServerError), err)
		http.Error(w, resp.Error(), http.StatusInternalServerError)
		log.Println(resp)
		return
	}
}

// verifyIntegrity checks that the criteria specified at the following link are satisfied:
// https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
func verifyIntegrity(idToken string) (*idtoken.Payload, error) {
	/*
		The ID token is properly signed by Google.
		Use Google's public keys (available in JWK or PEM format) to verify the token's signature.
		These keys are regularly rotated; examine the `Cache-Control` header in the response to determine when you should
		retrieve them again.
	*/
	idtPayload, err := idtoken.Validate(context.Background(), idToken, audience)
	if err != nil {
		return nil, fmt.Errorf("could not validate ID token: %w", err)
	}

	/*
		The value of `aud` in the ID token is equal to one of your app's client IDs.
		This check is necessary to prevent ID tokens issued to a malicious app being used to access data about the same
		user on your app's backend server.
	*/
	// This check should already have been made inside idtoken.Validate() above.
	if !strings.EqualFold(idtPayload.Audience, audience) {
		return nil, fmt.Errorf("token audience '%s' does not match this app's client ID", idtPayload.Audience)
	}

	/*
		The value of `iss` in the ID token is equal to `accounts.google.com` or `https://accounts.google.com`.
	*/
	if !strings.HasSuffix(idtPayload.Issuer, "accounts.google.com") {
		return nil, fmt.Errorf("token was issued by '%s' and not by Google Accounts", idtPayload.Issuer)
	}

	/*
		The expiry time (`exp`) of the ID token has not passed.
	*/
	tokenExpires := time.Unix(idtPayload.Expires, 0)
	if tokenExpires.Before(time.Now()) {
		return nil, fmt.Errorf("token already expired at '%s'", tokenExpires)
	}

	/*
		If you want to restrict access to only members of your G Suite domain, verify that the ID token has an `hd` claim
		that matches your G Suite domain name.
	*/

	// TODO: allowlist based on Google account ID

	// Everything checks out!
	return idtPayload, nil
}
