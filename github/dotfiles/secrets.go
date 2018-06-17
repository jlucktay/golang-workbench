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
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var tokenMap map[string]string

	if err := json.Unmarshal(fileContents, &tokenMap); err != nil {
		log.Fatal(err)
	}

	if _, ok := tokenMap["GitHubPersonalAccessToken"]; ok {
		token = tokenMap["GitHubPersonalAccessToken"]
	}

	return
}
