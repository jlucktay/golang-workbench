package main

import "fmt"

func main() {
	slice := []int{1, 2, 3, 4, 5}
	fmt.Println("slice:", slice)

	slice = slice[0:]
	fmt.Println("slice after zero chop:", slice)

	// chop off front 3 elements
	offset := 3

	// chop chop
	slice = slice[offset:]

	fmt.Println("slice after offset chop:", slice)

	// chop way too much
	slice = slice[100:]
}
