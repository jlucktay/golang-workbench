package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func main() {
	fmt.Println("Starting the application...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	opts := &options.ClientOptions{}
	opts = opts.SetConnectTimeout(10 * time.Second)
	opts = opts.SetHosts([]string{"mongodb://127.0.0.1:27017"})

	fmt.Printf("Connecting to Mongo at '%v'...\n", opts.Hosts)
	client, errConnect := mongo.Connect(ctx, opts)
	if errConnect != nil {
		log.Fatal(errConnect)
	}

	fmt.Println("Connected to Mongo!")

	fmt.Println("Pinging Mongo...")
	errPing := client.Ping(nil, nil)
	if errPing != nil {
		log.Fatal(errPing)
	}

	// db, errListDb := client.ListDatabaseNames(ctx, bson.D{}, &options.ListDatabasesOptions{})
	// if errListDb != nil {
	// 	log.Fatal(errListDb)
	// }
	// fmt.Printf("Mongo databases: %v\n", db)

	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
