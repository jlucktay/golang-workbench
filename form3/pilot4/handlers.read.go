package main

import (
	"encoding/json"
	"log"
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

		if payRead, errRead := a.storage.Read(id); errRead == nil {
			payBytes, errMarshal := json.Marshal(payRead)
			if errMarshal != nil {
				log.Fatal(errMarshal)
			}

			w.WriteHeader(http.StatusOK) // 200
			w.Header().Set("Content-Type", "application/json")
			if _, errWrite := w.Write(payBytes); errWrite != nil {
				log.Fatal(errWrite)
			}
			return
		}

		http.Error(w, (&NotFoundError{id}).Error(), http.StatusNotFound) // 404
	}
}
