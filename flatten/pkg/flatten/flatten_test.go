package flatten_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.jlucktay.dev/golang-workbench/flatten/pkg/flatten"
)

func TestFlatten(t *testing.T) {
	testCases := map[string]struct {
		input    any
		expected []int
	}{
		"Hello world - [[1,2,[3]],4] -> [1,2,3,4]": {
			input: []any{
				[]any{
					1,
					2,
					[]int{3},
				},
				4,
			},
			expected: []int{1, 2, 3, 4},
		},
		"Nested": {
			input: []any{
				[]any{
					1,
					[]int{2},
					[]any{
						3,
						[]int{4, 5},
						[]any{
							6,
							[]int{7, 8},
							9, 10,
						},
						11,
					},
					12,
				},
				13, 14,
			},
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		},
	}
	for desc, tC := range testCases {
		t.Run(desc, func(t *testing.T) {
			actual := flatten.Flatten(tC.input)

			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Got '%#v', want '%#v': mismatch (-want +got):\n%s", actual, tC.expected, diff)
			}
		})
	}
}
