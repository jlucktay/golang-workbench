package main

func (a *apiServer) setupRoutes() {
	// C
	a.router.POST("/payments", a.createPayments())
	a.router.POST("/payments/:id", a.createPaymentById())

	// R
	a.router.GET("/payments", a.readPayments())
	a.router.GET("/payments/:id", a.readPaymentById())

	// U
	a.router.PUT("/payments", a.updatePayments())
	a.router.PUT("/payments/:id", a.updatePaymentById())

	// D
	a.router.DELETE("/payments", a.deletePayments())
	a.router.DELETE("/payments/:id", a.deletePaymentById())
}
