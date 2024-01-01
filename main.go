package main

import (
	"log"
  "github.com/gin-gonic/gin"
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


	r.GET("/customers", GetAllCustomers)
	r.GET("/customers/:id", GetCustomer)


	r.Run(":8080")
}
