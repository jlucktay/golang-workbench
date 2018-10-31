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

func outputToJSON(jsonFilename string) {
	// Output the map of crawled URLs to a JSON file with current timestamp and
	// domain in its name. Range over the map, converting to strings and string
	// slices along the way, and copy into a slice of the custom type before
	// marshaling out to a JSON file.
	var cpSlice []CrawledPage

	for parent, children := range crawled.m {
		var cpChildren []string

		for _, child := range children {
			cpChildren = append(cpChildren, child.String())
		}

		cpSlice = append(cpSlice,
			CrawledPage{
				Parent:   parent.String(),
				Children: cpChildren,
			})
	}

	// Marshal the slice of custom types into JSON
	jsonBytes, errMarshal := json.MarshalIndent(cpSlice, "", "  ")
	if errMarshal != nil {
		Error.Printf("Error marshaling JSON: %v\n", errMarshal)
		return
	}

	// Emit the JSON to file
	errWrite := ioutil.WriteFile(jsonFilename, jsonBytes, 0644)
	if errWrite != nil {
		Error.Printf("Error writing to file '%s': %v\n", jsonFilename, errWrite)
		return
	}

	fmt.Println("Wrote page/link relationships to file:", jsonFilename)
}
