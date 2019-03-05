package main

import "net/http"

type apiServer struct {
	// db     *someDatabase

	// router *httprouter.Router
	router *http.ServeMux
}
