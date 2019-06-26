package flatten_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jlucktay/golang-workbench/flatten/pkg/flatten"
)

func TestFlatten(t *testing.T) {
	testCases := map[string]struct {
		input    interface{}
		expected []int
	}{
		"Hello world - [[1,2,[3]],4] -> [1,2,3,4]": {
			input: []interface{}{
				[]interface{}{
					1,
					2,
					[]int{3},
				},
				4,
			},
			expected: []int{1, 2, 3, 4},
		},
		"Nested": {
			input: []interface{}{
				[]interface{}{
					1,
					[]int{2},
					[]interface{}{
						3,
						[]int{4, 5},
						[]interface{}{
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
		tC := tC // pin!
		t.Run(desc, func(t *testing.T) {
			actual := flatten.Flatten(tC.input)

			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Got '%#v', want '%#v': mismatch (-want +got):\n%s", actual, tC.expected, diff)
			}
		})
	}
}
