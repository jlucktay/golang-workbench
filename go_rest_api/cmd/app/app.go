// Basic building blocks of the app.
package main

import (
	"fmt"
	"log"

	"go.jlucktay.dev/golang-workbench/go_rest_api/pkg/mongo"
	"go.jlucktay.dev/golang-workbench/go_rest_api/pkg/server"
)

// App has a server and a Mongo session.
type App struct {
	server  *server.Server
	session *mongo.Session
}

// Initialize sets up an App.
func (a *App) Initialize() {
	err := a.session.Open()
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	u := mongo.NewUserService(a.session.Copy())
	a.server = server.NewServer(u)
}

// Run will run an App.
func (a *App) Run() {
	fmt.Println("Run")
	defer a.session.Close()
	a.server.Start()
}
