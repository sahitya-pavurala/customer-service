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
        id := c.Param("id")

        if len(id)!=0{

        for _, customer := range customers {
                if customer.ID == id {
                        c.IndentedJSON(http.StatusOK, customer)
                        return
                }
        }
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "customer not found"})
        return
      }

      c.IndentedJSON(http.StatusOK, customers)

}

func main() {
  r := gin.Default()
  r.GET("/customers", getCustomers)
  r.GET("/customers/:id", getCustomers)

  r.Run(":8080")
}
