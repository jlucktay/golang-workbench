package main

import (
	"fmt"

	"go.jlucktay.dev/golang-workbench/flatten/pkg/flatten"
)

func main() {
	start := []interface{}{
		[]interface{}{
			1,
			2,
			[]int{3},
		},
		4,
	}

	finish := flatten.Flatten(start)

	fmt.Printf("Started with: %#v\n", start)
	fmt.Printf("Finished with: %#v\n", finish)
}
