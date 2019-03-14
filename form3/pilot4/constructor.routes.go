package main

func (a *apiServer) setupRoutes() {
	a.router.HandleMethodNotAllowed = true

	// C
	a.router.POST("/payments", a.createPayments())
	a.router.POST("/payments/:id", a.createPaymentById())

	// R
	a.router.GET("/payments", a.readPayments())
	a.router.GET("/payments/:id", a.readPaymentById())

	// U
	a.router.PUT("/payments/:id", a.updatePaymentById())

	// D
	a.router.DELETE("/payments/:id", a.deletePaymentById())
}
