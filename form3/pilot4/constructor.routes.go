package main

func (a *apiServer) setupRoutes() {
	// C
	a.router.POST("/payments", a.createPayment())
	a.router.POST("/payments/:id", a.createPaymentById())

	// R
	a.router.GET("/payments", a.readPayments())
	a.router.GET("/payments/:id", a.readPaymentById())

	// U
	a.router.PUT("/payments", a.updatePayment())

	// D
	a.router.DELETE("/payments", a.deletePayment())

	/*
		// Middleware
		s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))
	*/
}
