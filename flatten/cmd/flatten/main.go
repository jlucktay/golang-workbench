package main

import (
	"fmt"

	"github.com/jlucktay/golang-workbench/flatten/pkg/flatten"
)

func main() {
	input := []interface{}{
		[]interface{}{
			1,
			2,
			[]int{3},
		},
		4,
	}

	fmt.Printf("%#v\n", flatten.Flatten(input))
}
