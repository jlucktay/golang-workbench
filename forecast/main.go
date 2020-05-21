package main

import (
	"fmt"
	"log"

	"github.com/joefitzgerald/forecast"
	"go.jlucktay.dev/golang-workbench/secrets"
)

func main() {
	api := forecast.New(
		"https://api.forecastapp.com",
		"736067",
		secrets.ReadTokenFromSecrets("secrets.json"),
	)

	p, e := api.Projects()

	if e != nil {
		log.Fatal(e)
	}

	fmt.Println(p)
}
