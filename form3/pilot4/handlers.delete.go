package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *apiServer) deletePayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *apiServer) deletePaymentById() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
		w.WriteHeader(http.StatusNotImplemented)
	}
}
