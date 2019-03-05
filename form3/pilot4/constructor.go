package main

import (
	"github.com/julienschmidt/httprouter"
)

func newApiServer() (a *apiServer) {
	a = &apiServer{
		router: httprouter.New(),
	}
	a.setupRoutes()
	return
}
