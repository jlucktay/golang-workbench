package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) readPayments() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (a *apiServer) readPaymentById() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing

		if uuid.FromStringOrNil(p.ByName("id")) == uuid.Nil {
			http.Error(w, "invalid ID", http.StatusNotFound) // 404
			return
		}

		w.WriteHeader(http.StatusNotImplemented)
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

/*
// Handler-specific type(s)
func (s *server) handleSomething() http.HandlerFunc {
	type request struct {
		Name string
	}
	type response struct {
		Greeting string `json:"greeting"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		...
	}
}
*/

/*
// sync.Once to setup dependencies
func (s *server) handleTemplate(files string...) http.HandlerFunc {
	var (
		init sync.Once
		tpl  *template.Template
		err  error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		init.Do(func(){
			tpl, err = template.ParseFiles(files...)
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// use tpl
	}
}
*/
