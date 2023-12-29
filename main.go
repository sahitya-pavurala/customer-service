package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var db *bolt.DB

func getAllCustomers(c *gin.Context) {
	var data [][]byte
	var result map[string]interface{}
	results := []map[string]interface{}{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("customers"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		// Set the cursor to the last key in the bucket
		c := b.Cursor()
		lastKey, _ := c.Last()

		// Iterate and add to data
		for k, v := c.Seek(lastKey); k != nil; k, v = c.Prev() {
			data = append(data, v)
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for _, d := range data {
		json.Unmarshal(d, &result)
		results = append(results, result)
	}

	c.JSON(http.StatusOK, results)
}

func getCustomers(c *gin.Context) {
	id := c.Param("id")

	var result map[string]interface{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("customers"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("Customer not found")
		}

		return json.Unmarshal(data, &result)
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)

}

func main() {
	var err error
	db, err = bolt.Open("customer.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize database buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("customers"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("customers"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		// Seed data
		seedData := map[string]interface{}{
			"1": map[string]interface{}{
				"name":     "Luffy",
				"address":  "123 Main St",
				"accounts": []string{"acc1", "acc2"},
			},
			// Add more seed data as needed
		}

		for id, data := range seedData {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return err
			}

			err = b.Put([]byte(id), jsonData)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Seed data added to BoltDB successfully.")

	r := gin.Default()
	r.GET("/customers", getAllCustomers)
	r.GET("/customers/:id", getCustomers)

	r.Run(":8080")
}
