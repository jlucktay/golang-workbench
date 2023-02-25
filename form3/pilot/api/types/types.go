package types

import (
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	// db     *someDatabase
	router *httprouter.Router
	// email  EmailSender
}

func (s *Server) routes() {
	// s.router.HandleFunc()
}
