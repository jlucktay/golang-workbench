package pilot

import (
	"github.com/julienschmidt/httprouter"
)

type apiServer struct {
	/*
		db     *database
	*/
	router *httprouter.Router
}

func setup() {
	r := httprouter.New()
	r.GET("/payments", nil)
}
