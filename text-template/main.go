// How do text templates work?
package main

import (
	"os"
	"text/template"
)

// Inventory is a basic struct to aid the example below.
type Inventory struct {
	Material string
	Count    uint
}

func main() {
	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}\n")

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, sweaters)

	if err != nil {
		panic(err)
	}
}
