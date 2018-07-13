package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	ghpaToken = readTokenFromSecrets("./secrets.json")
)

type secrets struct {
	GitHubPersonalAccessToken string
}

func readTokenFromSecrets(path string) (token string) {
	fileContents, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var tokenMap map[string]string

	if umErr := json.Unmarshal(fileContents, &tokenMap); umErr != nil {
		log.Fatal(umErr)
	}

	if t, ok := tokenMap["GitHubPersonalAccessToken"]; ok {
		token = t
	}

	return
}
