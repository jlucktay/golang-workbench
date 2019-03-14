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
		w.WriteHeader(http.StatusNotImplemented) // 501
	}
}

func (a *apiServer) readPaymentById() httprouter.Handle {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// use thing

		if uuid.FromStringOrNil(p.ByName("id")) == uuid.Nil {
			http.Error(w, "Invalid ID.", http.StatusNotFound) // 404
			return
		}

		if p.ByName("id") == "29e1c453-8cc7-47b8-9c48-7e44b4f9ba26" {
			// Placeholder for valid route in the logic
			http.Error(w, "A payment with this ID does not exist.", http.StatusTeapot) // -> .StatusNotFound 404
		}

		w.WriteHeader(http.StatusNotImplemented) // 501
	}
}
