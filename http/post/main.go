package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type translateRequestBody struct {
	Q      []string `json:"q"`
	Target string   `json:"target"`
	Format string   `json:"format"`
}

type translateResponseBody struct {
	Data struct {
		Translations []struct {
			TranslatedText         string `json:"translatedText"`
			DetectedSourceLanguage string `json:"detectedSourceLanguage"`
		} `json:"translations"`
	} `json:"data"`
}

func main() {
	client := &http.Client{}

	reqBody := translateRequestBody{
		Q: []string{
			"你们好，我很高兴因为我在这里",
			"Kolik je hodin?",
			"予算内のすべて",
		},
		Target: "en",
		Format: "text",
	}

	reqBodyJSON, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatal(err)
	}

	reqVerb := "POST"
	reqURL := "https://translation.googleapis.com/language/translate/v2"

	req, err := http.NewRequest(reqVerb, reqURL, strings.NewReader(string(reqBodyJSON)))

	if err != nil {
		log.Fatal(err)
	}

	hdrAuth, err := exec.Command("gcloud", "auth", "print-access-token").Output()

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", strings.TrimSpace(string(hdrAuth))))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var respBody translateResponseBody

	err = json.Unmarshal(respBodyBytes, &respBody)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("respBody: '%+v'\n", respBody)
}
