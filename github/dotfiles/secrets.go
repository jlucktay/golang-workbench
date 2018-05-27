package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	ghpaToken = readTokenFromSecrets()
)

type secrets struct {
	GitHubPersonalAccessToken string
}

func readTokenFromSecrets() (token string) {
	fileContents, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatal(err)
	}

	var tokenMap map[string]string

	if err := json.Unmarshal(fileContents, &tokenMap); err != nil {
		log.Fatal(err)
	}

	token = tokenMap["GitHubPersonalAccessToken"]

	return
}
