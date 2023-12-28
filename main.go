package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

// customer struct to mock customer
type customer struct {
        ID     string  `json:"id"`
        Name  string  `json:"name"`
        AccountID     string  `json:"accountid"`
}

// customers slice to seed customer data.
var customers = []customer{
        {ID: "1", Name: "Luffy", AccountID: "123"},
        {ID: "2", Name: "Zoro", AccountID: "345"},
        {ID: "3", Name: "Sanji", AccountID: "678"},
}


func getCustomers(c *gin.Context) {
        c.IndentedJSON(http.StatusOK, customers)
}

func main() {
  r := gin.Default()
  r.GET("/customers", getCustomers)

  r.Run(":8080")
}
