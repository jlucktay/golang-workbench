package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	actorContent, actorErr := ioutil.ReadFile("actor.history.csv")
	if actorErr != nil {
		log.Fatal(actorErr)
	}

	fmt.Printf("Actor history: %s", actorContent)

	// taskContent, taskErr := ioutil.ReadFile("tasks.csv")
	// if taskErr != nil {
	// 	log.Fatal(taskErr)
	// }

	// fmt.Printf("Tasks: %s", taskContent)

	actorReader := csv.NewReader(bytes.NewReader(actorContent))

	actorRecords, actorReadErr := actorReader.ReadAll()
	if actorReadErr != nil {
		log.Fatal(actorReadErr)
	}

	fmt.Print(actorRecords)
}
