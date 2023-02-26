// Practising with JSON marshaling.
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type colorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {
	group := colorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	b, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
