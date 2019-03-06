package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestStatusCode(t *testing.T) {
	// Arrange
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Create a new payment",
			path:     "/payments",
			verb:     http.MethodPost,
			expected: http.StatusCreated,
		},
		{
			desc:     "Create a new payment on a pre-existing ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusConflict,
		},
		{
			desc:     "Create a new payment on a non-existent valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Create a new payment on an invalid ID",
			path:     "/payments/my-payment-id",
			verb:     http.MethodPost,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Read the entire collection of existing payments",
			path:     "/payments",
			verb:     http.MethodGet,
			expected: http.StatusOK,
		},
		{
			desc:     "Read a limited collection of existing payments",
			path:     "/payments?offset=2&limit=2",
			verb:     http.MethodGet,
			expected: http.StatusOK,
		},
		{
			desc:     "Read a single existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodGet,
			expected: http.StatusOK,
		},
		{
			desc:     "Read a non-existent payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodGet,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Read a non-existent payment at an invalid ID",
			path:     "/payments/my-payment-id",
			verb:     http.MethodGet,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Update all existing payments",
			path:     "/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Update an existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPut,
			expected: http.StatusNoContent,
		},
		{
			desc:     "Update a non-existent payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPut,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Update a non-existent payment at an invalid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPut,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Delete all existing payments",
			path:     "/payments",
			verb:     http.MethodDelete,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Delete an existing payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodDelete,
			expected: http.StatusOK,
		},
		{
			desc:     "Delete a non-existent payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Delete a non-existent payment at an invalid ID",
			path:     "/payments/my-payment-id",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound,
		},
	}

	srv := newApiServer()

	// Act & Assert
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			is := is.New(t)
			req, err := http.NewRequest(tC.verb, tC.path, nil)
			is.NoErr(err)
			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, req)
			is.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}
