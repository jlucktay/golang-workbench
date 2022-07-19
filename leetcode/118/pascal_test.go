package pascal

import (
	"reflect"
	"testing"

	"github.com/matryer/is"
)

func TestGenerate(t *testing.T) {
	testCases := map[string]struct {
		output [][]int
		input  int
	}{
		"Example 1": {
			input:  5,
			output: [][]int{{1}, {1, 1}, {1, 2, 1}, {1, 3, 3, 1}, {1, 4, 6, 4, 1}},
		},
		"Example 2": {
			input:  1,
			output: [][]int{{1}},
		},
	}
	for desc, tc := range testCases {
		t.Run(desc, func(t *testing.T) {
			is.New(t).True(reflect.DeepEqual(generate(tc.input), tc.output)) // expected != actual
		})
	}
}
