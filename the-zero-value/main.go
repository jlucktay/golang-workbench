// Go is reasonably nice and clever about default values for builtin variable types.
package main

import "fmt"

// T is an example struct to show off the zero value behaviour
type T struct {
	i    int
	f    float64
	next *T
}

func main() {
	t1 := new(T)
	fmt.Printf("%t\n", t1.i == 0)
	fmt.Printf("%t\n", t1.f == 0.0)
	fmt.Printf("%t\n", t1.next == nil)

	var t2 T
	fmt.Printf("%t\n", t2.i == 0)
	fmt.Printf("%t\n", t2.f == 0.0)
	fmt.Printf("%t\n", t2.next == nil)
}
