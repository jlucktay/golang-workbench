package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	actorContent, actorErr := os.ReadFile("actor.history.csv")
	if actorErr != nil {
		log.Fatal(actorErr)
	}

	fmt.Printf("Actor history: %s", actorContent)

	// taskContent, taskErr := os.ReadFile("tasks.csv")
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
