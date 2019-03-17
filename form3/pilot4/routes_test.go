package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	uuid "github.com/satori/go.uuid"
)

func TestCreateEmptyBody(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Create a new payment with an empty request body",
			path:     "/payments",
			verb:     http.MethodPost,
			expected: http.StatusBadRequest, // 400
		},
		{
			desc:     "Create a new payment on a pre-existing ID with an empty request body",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusConflict, // 409
		},
		{
			desc:     "Create a new payment on a non-existent valid ID with an empty request body",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Create a new payment on an invalid ID with an empty request body",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
	}

	srv := newApiServer(InMemory)

	for _, tC := range testCases {
		w := httptest.NewRecorder()

		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)

			req, err := http.NewRequest(tC.verb, tC.path, nil)
			i.NoErr(err)

			srv.router.ServeHTTP(w, req)
			i.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}

		},
		{
			desc:     "Create a new payment on a pre-existing ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusConflict, // 409
		},
		{
			desc:     "Create a new payment on a non-existent valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Create a new payment on an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPost,
			expected: http.StatusNotFound, // 404
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Fatalf("not yet implemented")
		})
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Read the entire collection of existing payments",
			path:     "/payments",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		{
			desc:     "Read a limited collection of existing payments",
			path:     "/payments?offset=2&limit=2",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		{
			desc:     "Read a single existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodGet,
			expected: http.StatusOK, // 200
		},
		{
			desc:     "Read a non-existent payment at a valid ID",
			path:     fmt.Sprintf("/payments/%s", uuid.Must(uuid.NewV4())),
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Read a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodGet,
			expected: http.StatusNotFound, // 404
		},
	}

	srv := newApiServer(InMemory)

	for _, tC := range testCases {
		w := httptest.NewRecorder()

		t.Run(tC.desc, func(t *testing.T) {
			i := is.New(t)

			req, err := http.NewRequest(tC.verb, tC.path, nil)
			i.NoErr(err)

			srv.router.ServeHTTP(w, req)
			i.Equal(tC.expected, w.Result().StatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Update all existing payments",
			path:     "/payments",
			verb:     http.MethodPut,
			expected: http.StatusMethodNotAllowed, // 405
		},
		{
			desc:     "Update an existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPut,
			expected: http.StatusNoContent, // 204; update is OK, but response has no body/content
		},
		{
			desc:     "Update a non-existent payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Update a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodPut,
			expected: http.StatusNotFound, // 404
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Fatalf("not yet implemented")
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Delete all existing payments",
			path:     "/payments",
			verb:     http.MethodDelete,
			expected: http.StatusMethodNotAllowed, // 405
		},
		{
			desc:     "Delete an existing payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodDelete,
			expected: http.StatusOK, // 200
		},
		{
			desc:     "Delete a non-existent payment at a valid ID",
			path:     "/payments/1234-5678-abcd",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
		{
			desc:     "Delete a non-existent payment at an invalid ID",
			path:     "/payments/not-a-valid-v4-uuid",
			verb:     http.MethodDelete,
			expected: http.StatusNotFound, // 404
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Fatalf("not yet implemented")
		})
	}
}
