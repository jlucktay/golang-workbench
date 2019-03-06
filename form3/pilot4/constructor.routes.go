package main

func (a *apiServer) setupRoutes() {
	// C
	a.router.POST("/payments", a.createPayment())

	// R
	a.router.GET("/payments", a.readAllPayments())
	a.router.GET("/payments/:id", a.readSinglePayment())

	// U
	a.router.PUT("/payments", a.updatePayment())

	// D
	a.router.DELETE("/payments", a.deletePayment())

	/*
		// Middleware
		s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
	*/
}
