// Create a struct that describes only part of a JSON object and unmarshal into it.
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	jsonBlob := []byte(`{
		"ID": 1,
		"Name": "Reds",
		"Colors": [
			"Crimson",
			"Red",
			"Ruby",
			"Maroon"
		]
	}`)

	type partialJSON struct {
		Name string
	}

	var j partialJSON

	if umErr := json.Unmarshal(jsonBlob, &j); umErr != nil {
		log.Fatal(umErr)
	}

	fmt.Printf("%+v\n", j)
}
