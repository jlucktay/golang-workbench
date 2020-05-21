// Package server defines internal behaviour.
package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	root "go.jlucktay.dev/golang-workbench/go_rest_api/pkg"
)

// Server holds a router.
type Server struct {
	router *mux.Router
}

// NewServer returns an Server initialised with the given UserService.
func NewServer(u root.UserService) *Server {
	s := Server{router: mux.NewRouter()}

	s.router.Handle("/user", NewUserRouter(u, s.getSubrouter("/user")))

	return &s
}

// Start will start the Server listening.
func (s *Server) Start() {
	log.Println("Listening on port 1337...")
	if err := http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
}

func (s *Server) getSubrouter(path string) *mux.Router {
	return s.router.PathPrefix(path).Subrouter()
}
