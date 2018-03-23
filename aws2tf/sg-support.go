package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"./aws2tf"
)

func main() {
	raw, err := ioutil.ReadFile("sg-support.json")

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
