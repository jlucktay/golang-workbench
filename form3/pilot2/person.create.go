package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	var person Person

	errDecode := json.NewDecoder(request.Body).Decode(&person)
	if errDecode != nil {
		log.Fatal(errDecode)
	}

	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)

	errEncode := json.NewEncoder(response).Encode(result)
	if errEncode != nil {
		log.Fatal(errEncode)
	}
}
