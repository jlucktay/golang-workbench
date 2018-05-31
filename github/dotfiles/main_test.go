package main

import (
	"fmt"
	"testing"
)

func TestMatchStringCompiled(t *testing.T) {
	testCases := []struct {
		needle   string
		haystack string
		expected bool
	}{
		{
			needle:   "needle",
			haystack: "haystack",
			expected: false,
		},
		{
			needle:   "github",
			haystack: "https://api.github.com/repos/jlucktay/adventofcode/languages",
			expected: true,
		},
		{
			needle:   `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`,
			haystack: "https://github.com/atomantic/dotfiles",
			expected: true,
		},
		{
			needle:   `^https://github.com/[A-Za-z0-9_\-\.]+/[A-Za-z0-9_\-\.]+$`,
			haystack: "https://www.bro.org/",
			expected: false,
		},
	}

	for _, tC := range testCases {
		desc := fmt.Sprintf("'%s' in '%s'", tC.needle, tC.haystack)

		t.Run(desc, func(t *testing.T) {
			if matchStringCompiled(tC.needle, tC.haystack) != tC.expected {
				t.Fatalf("'%s' failed", desc)
			}
		})
	}
}
