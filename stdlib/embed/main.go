package main

import (
	_ "embed"
)

//go:embed main.go
var s string

func main() {
	println(s)
}
