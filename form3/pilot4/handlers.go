package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *apiServer) handleIndex() http.HandlerFunc {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
	}
}

func (a *apiServer) handleGetAllPayments() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
	}
}

func (a *apiServer) handleGetSinglePayment() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
	}
}

/*
// Middleware handler
func (a *apiServer) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !currentUser(r).IsAdmin {
			http.NotFound(w, r)
			return
		}
		h(w, r)
	}
}
*/
