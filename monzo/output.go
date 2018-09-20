package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// CrawledPage is a custom type, for holding parent/child page relationships
type CrawledPage struct {
	Parent   string
	Children []string
}

func outputToJSON() {
	// Output the map of crawled URLs to a JSON file with current timestamp and domain in its name
	// Range over the map, converting to string/string slices along the way, and copy into a slice of the custom type
	cpSlice := make([]CrawledPage, 0)

	for a, b := range crawled.m {
		cpChildren := make([]string, 0)

		for _, c := range b {
			cpChildren = append(cpChildren, c.String())
		}

		cpSlice = append(cpSlice,
			CrawledPage{
				Parent:   a.String(),
				Children: cpChildren,
			})
	}

	// Marshal the slice of custom types into JSON
	b, errMarshal := json.MarshalIndent(cpSlice, "", "  ")
	if errMarshal != nil {
		fmt.Println("error:", errMarshal)
	}

	// Emit the JSON to file
	jsonFilename := timestamp + "." + flagURL + ".json"
	errWrite := ioutil.WriteFile(jsonFilename, b, 0644)
	if errWrite != nil {
		fmt.Println("error:", errWrite)
	}
}
