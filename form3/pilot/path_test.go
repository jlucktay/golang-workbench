package main

import (
	"testing"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc     string
		path     string
		verb     string
		expected int
	}{
		{
			desc:     "Create a new payment",
			path:     "/payments",
			verb:     "POST",
			expected: 0,
		},
		{
			desc:     "Read a single existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     "GET",
			expected: 0,
		},
		{
			desc:     "Read a limited collection of existing payments",
			path:     "/payments?offset=2&limit=2",
			verb:     "GET",
			expected: 0,
		},
		{
			desc:     "Update an existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     "PUT",
			expected: 0,
		},
		{
			desc:     "Delete an existing payment",
			path:     "/payments/1234-5678-abcd",
			verb:     "DELETE",
			expected: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
