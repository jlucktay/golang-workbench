package main

import (
	"encoding/json"
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

	for parent, children := range crawled.m {
		cpChildren := make([]string, len(children))

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
		errorLog.Printf("Error marshaling JSON: %v\n", errMarshal)
	}

	// Emit the JSON to file
	jsonFilename := timestamp + "." + flagURL + ".json"
	errWrite := ioutil.WriteFile(jsonFilename, jsonBytes, 0644)
	if errWrite != nil {
		errorLog.Printf("Error writing to file '%s': %v\n", jsonFilename, errWrite)
	}
}
