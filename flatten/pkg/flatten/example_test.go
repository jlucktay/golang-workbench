package flatten_test

import (
	"fmt"

	"github.com/jlucktay/golang-workbench/flatten/pkg/flatten"
)

func ExampleFlatten() {
	fmt.Printf("%#v\n", flatten.Flatten([]int{1, 2, 3}))
	// Output: []int{1, 2, 3}
}

func ExampleFlatten_complex() {
	complex := []interface{}{
		[]interface{}{
			1,
			[]int{2, 3},
		},
		4,
		[]int{5, 6, 7},
		8,
	}
	fmt.Printf("%#v\n", flatten.Flatten(complex))
	// Output: []int{1, 2, 3, 4, 5, 6, 7, 8}
}
