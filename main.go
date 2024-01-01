package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := InitializeDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = SeedData()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// Public route for api user registration
	r.POST("/register", RegisterHandler)

	// Public route for login
	r.POST("/login", LoginHandler)

	// Public route for monitoring
	r.GET("/health", HealthHandler)

	// Secured route for a customers endpoint
	private := r.Group("/api")
	private.Use(AuthMiddleware)
	{
		private.GET("/customers", GetAllCustomers)
		private.GET("/customers/:id", GetCustomerUsingId)
	}

	r.Run(":8080")
}
