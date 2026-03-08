// Makes use of my homegrown AWS CLI JSON parser.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go.jlucktay.dev/golang-workbench/aws2tf/aws2tf"
)

func main() {
	raw, err := os.ReadFile("sg-support.json")
	if err != nil {
		panic(err)
	}

	var sgsupport aws2tf.SGFile

	err = json.Unmarshal(raw, &sgsupport)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", sgsupport)
}
