package kipa_test

import (
	"testing"

	kipa "go.jlucktay.dev/golang-workbench/leetcode/629"
)

func TestKInversePairst(t *testing.T) {
	testCases := map[string]struct {
		n, k     int
		expected int
	}{
		"Example 1": {
			n: 3, k: 0,
			expected: 1,
		},
		"Example 2": {
			n: 3, k: 1,
			expected: 2,
		},
		"https://www.youtube.com/watch?v=dD7jopXly08 - n=3,k=0": {
			n: 3, k: 0,
			expected: 1,
		},
		"https://www.youtube.com/watch?v=dD7jopXly08 - n=3,k=1": {
			n: 3, k: 1,
			expected: 2,
		},
		"https://www.youtube.com/watch?v=dD7jopXly08 - n=3,k=2": {
			n: 3, k: 2,
			expected: 2,
		},
		"https://www.youtube.com/watch?v=dD7jopXly08 - n=3,k=3": {
			n: 3, k: 3,
			expected: 1,
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			if result := kipa.KInversePairs(tc.n, tc.k); result != tc.expected {
				t.Fatalf("%d != %d", result, tc.expected)
			}
		})
	}
}
