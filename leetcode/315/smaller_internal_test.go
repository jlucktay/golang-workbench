package smaller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCountSmaller(t *testing.T) {
	testCases := map[string]struct {
		input          []int
		expectedOutput []int
	}{
		"Example 1": {
			input:          []int{5, 2, 6, 1},
			expectedOutput: []int{2, 1, 1, 0},
			// Explanation:
			// To the right of 5 there are 2 smaller elements (2 and 1).
			// To the right of 2 there is only 1 smaller element (1).
			// To the right of 6 there is 1 smaller element (1).
			// To the right of 1 there is 0 smaller element.
		},
		"Example 2": {
			input:          []int{-1},
			expectedOutput: []int{0},
		},
		"Example 3": {
			input:          []int{-1, -1},
			expectedOutput: []int{0, 0},
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			if diff := cmp.Diff(tc.expectedOutput, countSmaller(tc.input)); diff != "" {
				t.Errorf("countSmaller(%d) mismatch (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}
