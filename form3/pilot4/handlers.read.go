package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

func (a *apiServer) readPayments() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented) // 501
	}
}

func (a *apiServer) readPaymentById() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.FromStringOrNil(p.ByName("id"))

		if id == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		// Placeholder for valid route in the logic:
		// -> Read a non-existent payment at a valid ID
		if id.String() == "29e1c453-8cc7-47b8-9c48-7e44b4f9ba26" {
			http.Error(w, (&NotFoundError{id}).Error(), http.StatusTeapot) // -> .StatusNotFound 404
		}

		w.WriteHeader(http.StatusNotImplemented) // 501
	}
}
