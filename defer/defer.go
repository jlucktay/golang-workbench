package main

import "fmt"

func main() {
	var dp deferredPrinter

	dp = &one{}
	defer dp.deferredPrint("printing from one")

	dp = &two{}
	defer dp.deferredPrint("printing from two")
}

type deferredPrinter interface {
	deferredPrint(string)
}

type one struct{}

func (o *one) deferredPrint(s string) {
	fmt.Println(s)
}

type two struct{}

func (t *two) deferredPrint(s string) {
	fmt.Println(s)
	fmt.Println(s)
}
