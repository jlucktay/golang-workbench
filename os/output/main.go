package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("gcloud", "auth", "print-access-token").Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Output: '%s'\n", strings.TrimSpace(string(out)))
}

// gcloud auth print-access-token
