package main

import (
	"log"

	"github.com/jlucktay/golang-workbench/go_rest_api/pkg/crypto"
	"github.com/jlucktay/golang-workbench/go_rest_api/pkg/mongo"
	"github.com/jlucktay/golang-workbench/go_rest_api/pkg/server"
)

func main() {
	ms, err := mongo.NewSession("127.0.0.1:27017")
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	defer ms.Close()

	h := crypto.Hash{}
	u := mongo.NewUserService(ms.Copy(), "go_web_server", "user", &h)
	s := server.NewServer(u)

	s.Start()
}
