package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	fmt.Println("GetPeopleEndpoint - start")
	response.Header().Set("content-type", "application/json")

	var people []Person

	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, errWrite := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if errWrite != nil {
			log.Fatal(errWrite)
		}
		return
	}
	defer cursor.Close(ctx)

	fmt.Println("GetPeopleEndpoint - before cursor.Next()")
	for cursor.Next(ctx) {
		fmt.Println("GetPeopleEndpoint - cursor.Next()")
		var person Person
		errDecode := cursor.Decode(&person)
		if errDecode != nil {
			log.Fatal(errDecode)
		}
		people = append(people, person)
	}
	fmt.Println("GetPeopleEndpoint - after cursor.Next()")

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, errWrite := response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		if errWrite != nil {
			log.Fatal(errWrite)
		}
		return
	}

	errEncode := json.NewEncoder(response).Encode(people)
	if errEncode != nil {
		log.Fatal(errEncode)
	}
	fmt.Println("GetPeopleEndpoint - finish")
}
