package main

import (
	"fmt"
	"testing"
)

func TestFilterForGoRepos(t *testing.T) {
	testCases := []struct {
		languages map[string]int
		expected  bool
	}{
		{
			languages: map[string]int{"Go": 24031},
			expected:  true,
		},
		{
			languages: map[string]int{"Perl": 19111, "Shell": 8593, "Perl 6": 2945, "Makefile": 1355, "Ruby": 935, "Go": 401},
			expected:  false,
		},
		{
			languages: map[string]int{"Shell": 12464},
			expected:  false,
		},
		{
			languages: map[string]int{"Shell": 6455},
			expected:  false,
		},
		{
			languages: map[string]int{"Go": 24369, "Perl": 2497, "Makefile": 1582},
			expected:  true,
		},
		{
			languages: map[string]int{"Haskell": 193968, "Shell": 3663},
			expected:  false,
		},
		{
			languages: map[string]int{"Shell": 23490, "Ruby": 1444},
			expected:  false,
		},
		{
			languages: map[string]int{"Vim script": 31878, "Shell": 6547},
			expected:  false,
		},
		{
			languages: map[string]int{"Shell": 8131, "Batchfile": 5516},
			expected:  false,
		},
		{
			languages: map[string]int{"CSS": 7643, "Ruby": 6012, "HTML": 2970, "JavaScript": 537, "Shell": 218},
			expected:  false,
		},
	}

	for _, tC := range testCases {
		desc := fmt.Sprint(tC.languages)

		t.Run(desc, func(t *testing.T) {
			if filterForGoRepos(tC.languages) != tC.expected {
				t.Fatalf("'%s' failed", desc)
			}
		})
	}
}
