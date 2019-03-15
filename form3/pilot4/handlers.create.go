package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) createPayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (a *apiServer) createPaymentById() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		// Placeholder for valid route in the logic:
		// -> Create a new payment on a pre-existing ID
		if id.String() == "b2e3ccaa-ac37-45e0-b889-1e6acadf31c8" {
			http.Error(w, (&AlreadyExistsError{id}).Error(), http.StatusTeapot) // -> .StatusConflict 409
		}

		http.Error(w, `Cannot specify an ID for payment creation.
One will be generated for you.`, http.StatusNotFound) // 404
	}
}
