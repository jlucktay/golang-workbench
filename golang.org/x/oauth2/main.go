// Package main uses golang.org/x/oauth2/github to simplify the workflow of OAuth 2.0 with GitHub.
package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// These should be taken from your [GitHub application settings].
//
// [GitHub application settings]: https://github.com/settings/developers
var (
	GithubClientID     = os.Getenv("GITHUB_CLIENT_ID")
	GithubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
)

const (
	stateSize     = 16
	serverTimeout = 5 * time.Second
)

func main() {
	if len(GithubClientID) == 0 || len(GithubClientSecret) == 0 {
		slog.Error("set GITHUB_CLIENT_* env vars")

		return
	}

	addr := "localhost:8080"
	callbackPath := "/github/callback/"

	// Note: GitHub auth doesn't support PKCE verification, so we don't need [oauth2.GenerateVerifier].
	// See the [Proof Key for Code Exchange by OAuth Public Clients RFC] for more details.
	//
	// [Proof Key for Code Exchange by OAuth Public Clients RFC]: https://www.rfc-editor.org/rfc/rfc7636.html
	conf := &oauth2.Config{
		ClientID:     GithubClientID,
		ClientSecret: GithubClientSecret,
		Scopes:       []string{},
		Endpoint:     github.Endpoint,
	}

	flow := &loginFlow{
		config: conf,
	}

	http.HandleFunc("/", flow.rootHandler)
	http.HandleFunc("/login/", flow.githubLoginHandler)
	http.HandleFunc(callbackPath, flow.githubCallbackHandler)

	srv := http.Server{
		Addr:         addr,
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
	}

	slog.Info("listening", slog.String("address", srv.Addr))
	slog.Error(srv.ListenAndServe().Error())
}

const rootHTML = `<h1>My web app</h1>
<p>Using the x/oauth2 package</p>
<p>You can log into this app with your GitHub credentials:</p>
<p><a href="/login/">Log in with GitHub</a></p>
`

type loginFlow struct {
	config *oauth2.Config
}

func (lf *loginFlow) rootHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, rootHTML)
}

func (lf *loginFlow) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random state for CSRF protection and set it in a cookie.
	state, err := randString(stateSize)
	if err != nil {
		panic(err)
	}

	c := &http.Cookie{
		Name:     "state",
		Value:    state,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)

	redirectURL := lf.config.AuthCodeURL(state)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (lf *loginFlow) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)

		return
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)

		return
	}

	code := r.URL.Query().Get("code")

	tok, err := lf.config.Exchange(r.Context(), code)
	if err != nil {
		panic(err)
	}

	// This client will have a bearer token to access the GitHub API on the user's behalf.
	client := lf.config.Client(r.Context(), tok)

	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	respbody, _ := io.ReadAll(resp.Body)
	userInfo := string(respbody)

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, userInfo)
}

// randString generates a random string of length n and returns its base64-encoded version.
func randString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("reading from crypto/rand: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
