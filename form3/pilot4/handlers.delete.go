package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *apiServer) deletePaymentById() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
