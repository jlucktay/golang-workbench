package main

import (
	"fmt"
	"log"

	"github.com/jlucktay/golang-workbench/secrets"
	"github.com/joefitzgerald/forecast"
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
