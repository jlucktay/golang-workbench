package main

func (a *apiServer) routes() {
	// C
	// s.router.POST("/payments")

	// R
	a.router.GET("/payments", a.handleGetAllPayments())
	a.router.GET("/payments/:id", a.handleGetSinglePayment())

	// U
	// s.router.PUT("/payments", nil)

	// D
	// s.router.DELETE("/payments", nil)

	/*
		// Middleware
		s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
	*/
}
