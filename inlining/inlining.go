// See the accompanying README.
// The Go compiler is pretty cool/smart!
package main

import (
	"fmt"
)

func main() {
	a := 20
	b := 7

	c := add(a, b)

	fmt.Println(c)
}

func add(x, y int) int {
	return x + y
}
