package main

import "net/http"

func newApiServer() *apiServer {
	a := &apiServer{
		router: http.NewServeMux(),
	}
	a.routes()
	return a
}
