package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 2 {
		panic("provide exactly one arg to password-ify")
	}

	plaintext := []byte(os.Args[1])

	pass, err := bcrypt.GenerateFromPassword(plaintext, 10)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(pass, plaintext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n%s\n", pass, time.Now().Format(time.RFC3339))
}
