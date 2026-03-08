// The 'defer' keyword will defer execution until the end of that block.
package main

import "fmt"

func deferHello() {
	defer fmt.Println("world")

	fmt.Println("hello")
}

func deferCount() {
	fmt.Println("counting")

	for i := range 10 {
		defer fmt.Println(i)
	}

	fmt.Println("done")
}
