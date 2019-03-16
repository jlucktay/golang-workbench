package main

import (
	"log"

	"github.com/julienschmidt/httprouter"
)

func newApiServer(st StorageType) (a *apiServer) {
	a = &apiServer{
		router: httprouter.New(),
	}
	a.setupRoutes()

	switch st {
	case InMemory:
		a.storage = &inMemoryStorage{}
	case Mongo:
		panic("Mongo storage is not yet implemented")
	}

	if errStorageInit := a.storage.Init(); errStorageInit != nil {
		log.Fatal(errStorageInit)
	}

	return
}
