package main

import "fmt"

func deferHello() {
	defer fmt.Println("world")

	fmt.Println("hello")
}

func deferCount() {
	fmt.Println("counting")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("done")
}
