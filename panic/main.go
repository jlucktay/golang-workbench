package main

import "fmt"

type MakeMePanic struct {
	field string
}

func main() {
	one()
}

func one() {
	two()
}

func two() {
	three()
}

func three() {
	var mmp *MakeMePanic

	// Attempt to print the string field without creating the surrounding struct
	fmt.Printf("%v\n", mmp.field)
}
