package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	envKey := "VERSION"

	if ver, ok := os.LookupEnv(envKey); ok {
		fmt.Printf("Hello Docker %s!\n", ver)
	} else {
		log.Fatalf("'$%s' not found in environment\n", envKey)
	}
}
