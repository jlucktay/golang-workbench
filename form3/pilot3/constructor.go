package main

import "net/http"

func newApiServer() (a *apiServer) {
	a = &apiServer{
		router: http.NewServeMux(),
	}
	a.routes()
	return
}
