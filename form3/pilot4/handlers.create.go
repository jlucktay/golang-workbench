package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) createPayments() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (a *apiServer) createPaymentById() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if uuid.FromStringOrNil(p.ByName("id")) == uuid.Nil {
			http.Error(w, "invalid ID", http.StatusNotFound) // 404
			return
		}
		http.Error(w, "cannot specify an ID for payment creation", http.StatusConflict) // 409
	}
}
