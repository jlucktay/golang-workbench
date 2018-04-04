// Package server defines internal behaviour.
package server

import (
	"encoding/json"
	"net/http"
)

// Error explains what went wrong.
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}

// JSON marshals a JSON payload and writes it out to the response, with some headers.
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// JSONWithCookie marshals a JSON payload, sets a cookie, and writes out to the response, with some headers.
func JSONWithCookie(w http.ResponseWriter, code int, payload interface{}, cookie http.Cookie) {
	response, _ := json.Marshal(payload)
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
