package main

import "fmt"

func main() {
	fmt.Print("Enter text: ")

	var input string

	if _, err := fmt.Scanln(&input); err != nil {
		panic(err)
	}

	fmt.Println(input)
}
