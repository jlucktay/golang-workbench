package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/dghubble/sessions"
	"github.com/yosssi/gohtml"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

const (
	sessionName   = "example-google-app"
	sessionSecret = "example cookie signing secret"

	sessionEmail    = "googleEmail"
	sessionPicture  = "googlePicture"
	sessionUserKey  = "googleID"
	sessionUsername = "googleName"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// Config configures the main ServeMux.
type Config struct {
	ClientID, ClientSecret string
}

// New returns a new ServeMux with app routes.
func New(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", profileHandler)
	mux.HandleFunc("/logout", logoutHandler)

	// 1. Register Login and Callback handlers
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://localhost:8080/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig

	mux.Handle(
		"/google/login",
		google.StateHandler(
			stateConfig,
			google.LoginHandler(oauth2Config, nil),
		),
	)

	mux.Handle(
		"/google/callback",
		google.StateHandler(
			stateConfig,
			google.CallbackHandler(oauth2Config, issueSession(), nil),
		),
	)

	return mux
}

// issueSession issues a cookie session after successful Google login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("[issueSession] %s %s\n", req.Method, req.URL)

		googleUser, err := google.UserFromContext(req.Context())
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get Google user info from context: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionEmail] = googleUser.Email
		session.Values[sessionPicture] = googleUser.Picture
		session.Values[sessionUserKey] = googleUser.Id
		session.Values[sessionUsername] = googleUser.Name

		if err := session.Save(w); err != nil {
			fmt.Fprintf(os.Stderr, "could not save session: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/profile", http.StatusFound)
	}

	return http.HandlerFunc(fn)
}

// profileHandler shows a personal profile or a login button (unauthenticated).
func profileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("[profileHandler] %s %s\n", req.Method, req.URL)

	session, err := sessionStore.Get(req, sessionName)
	if err != nil {
		// welcome with login button
		page, err := ioutil.ReadFile("home.html")
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read HTML file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(page))

		return
	}

	// authenticated profile
	tpl := template.Must(template.New("profile.gohtml").ParseFiles("profile.gohtml"))
	data := struct {
		Items   map[string]interface{}
		Picture string
	}{
		Items: session.Values,
	}

	if picture, hasPicture := session.Values[sessionPicture]; hasPicture {
		if picStr, ok := picture.(string); ok {
			data.Picture = picStr
		}
	}

	b := &bytes.Buffer{}
	if err := tpl.Execute(b, data); err != nil {
		fmt.Fprintf(os.Stderr, "could not execute template into buffer: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fBytes := gohtml.Format(b.String())
	if _, err := w.Write([]byte(fBytes)); err != nil {
		fmt.Fprintf(os.Stderr, "could not write formatted bytes to ResponseWriter: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("[logoutHandler] %s %s\n", req.Method, req.URL)

	if req.Method == http.MethodPost {
		sessionStore.Destroy(w, sessionName)
	}

	http.Redirect(w, req, "/", http.StatusFound)
}

// main creates and starts a Server listening.
func main() {
	config := &Config{
		// read credentials from environment variables if available
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	// allow consumer credential flags to override config fields
	clientID := flag.String("client-id", "", "Google Client ID")
	clientSecret := flag.String("client-secret", "", "Google Client Secret")

	// default address to localhost for development
	address := flag.String("server-address", "localhost:8080", "Server address to listen on")
	flag.Parse()

	if *clientID != "" {
		config.ClientID = *clientID
	}

	if *clientSecret != "" {
		config.ClientSecret = *clientSecret
	}

	if config.ClientID == "" {
		log.Fatal("Missing Google Client ID")
	}

	if config.ClientSecret == "" {
		log.Fatal("Missing Google Client Secret")
	}

	log.Printf("Server listening on '%s'...\n", *address)
	log.Fatal(http.ListenAndServe(*address, New(config)))
}
