package main

import "fmt"

func main() {
	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	fmt.Println(input)
}
