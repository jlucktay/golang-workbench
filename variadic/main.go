package main

import "fmt"

func main() {
	variadicInputs("hello", "world")              // valid; multiple single <T>s
	variadicInputs([]string{"hello", "world"}...) // valid; single []<T> with ... expansion
	// variadicInputs("hello", []string{"world"}...) // not valid, can't mix
	// variadicInputs([]string{"hello"}..., "world") // not valid; can't mix
}

func variadicInputs(input ...string) {
	for a, b := range input {
		fmt.Printf("%v %v\n", a, b)
	}
}
