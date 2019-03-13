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
			path:     "/payments/b2e3ccaa-ac37-45e0-b889-1e6acadf31c8",
			verb:     http.MethodPost,
			expected: http.StatusConflict,
		},
		{
			desc:     "Create a new payment on a non-existent valid ID",
			path:     "/payments/60c4feb1-bf67-488a-8d04-627bac487c05",
			verb:     http.MethodPost,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Create a new payment on an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
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
			path:     "/payments/29e1c453-8cc7-47b8-9c48-7e44b4f9ba26",
			verb:     http.MethodGet,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Read a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
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
			path:     "/payments/ac5f6fcd-8e69-4949-ad93-d15c51991802",
			verb:     http.MethodPut,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Update a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
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
			path:     "/payments/a300eb47-efe0-44b0-b729-bed75123bf3a",
			verb:     http.MethodDelete,
			expected: http.StatusOK,
		},
		{
			desc:     "Delete a non-existent payment at a valid ID",
			path:     "/payments/943c4811-f66a-4cb1-8d5b-3ed7da0ce934",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound,
		},
		{
			desc:     "Delete a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound,
		},
	}

	srv := newApiServer()

	// Act & Assert
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)
			req, err := http.NewRequest(tC.verb, tC.path, nil)
			i.NoErr(err)
			w := httptest.NewRecorder()
			srv.router.ServeHTTP(w, req)
			i.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}
