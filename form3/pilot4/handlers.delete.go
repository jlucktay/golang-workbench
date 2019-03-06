package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *apiServer) deletePayment() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing
		w.WriteHeader(http.StatusNotImplemented)
	}
}
