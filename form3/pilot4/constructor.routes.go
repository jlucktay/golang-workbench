package main

func (a *apiServer) setupRoutes() {
	// C
	// s.router.POST("/payments")

	// R
	a.router.GET("/payments", a.getAllPayments())
	a.router.GET("/payments/:id", a.getSinglePayment())

	// U
	// s.router.PUT("/payments", nil)

	// D
	// s.router.DELETE("/payments", nil)

	/*
		// Middleware
		s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
	*/
}
