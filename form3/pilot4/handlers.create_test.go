package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestCreateNewPayment(t *testing.T) {
	i := is.New(t)
	a := newApiServer(InMemory)
	var w *httptest.ResponseRecorder

	// Construct a HTTP request which creates a payment
	p := Payment{Amount: decimal.NewFromFloat(123.45)}
	j, errMarshal := json.Marshal(p)
	i.NoErr(errMarshal)
	reqBody := bytes.NewBuffer(j)
	reqCreate, errCreate := http.NewRequest(http.MethodPost, "/payments", reqBody)
	i.NoErr(errCreate)
	reqCreate.Header.Set("Content-Type", "application/json")

	// Send it, and gather the ID of the new payment
	w = httptest.NewRecorder()
	a.router.ServeHTTP(w, reqCreate)
	i.Equal(http.StatusCreated, w.Result().StatusCode)

	// Make sure the response had a Location header pointing at the new payment
	loc := w.Result().Header.Get("Location")
	r := regexp.MustCompile("^/payments/([0-9a-f-]{36})$")
	i.True(r.MatchString(loc))
	newId := r.FindStringSubmatch(loc)[1]

	// Construct another HTTP request to read the new payment
	reqRead, errRead := http.NewRequest(http.MethodGet, "/payments/"+newId, nil)
	i.NoErr(errRead)

	// Read the payment under the given ID, as it should exist now
	w = httptest.NewRecorder()
	a.router.ServeHTTP(w, reqRead)
	i.Equal(http.StatusOK, w.Result().StatusCode)
}
