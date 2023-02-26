package main

import "net/http"

func (a *apiServer) routes() {
	// C
	// s.router.POST("/payments")

	// R
	a.router.HandleFunc("/payments", a.GET(a.handleGetSinglePayment()))
	// s.router.GET("/payments", nil)

	// U
	// s.router.PUT("/payments", nil)

	// D
	// s.router.DELETE("/payments", nil)

	//.HandleFunc("/api/", s.handleAPI())
	// s.router.HandleFunc("/about", s.handleAbout())
	// s.router.HandleFunc("/", s.handleIndex())
}

/*
func (s *server) handleSomething() http.HandlerFunc {
    thing := prepareThing()
    return func(w http.ResponseWriter, r *http.Request) {
        // use thing
    }
}
*/

func (a *apiServer) handleIndex() http.HandlerFunc {
	// thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) {
		// use thing
	}
}
