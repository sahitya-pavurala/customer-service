package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllCustomers(c *gin.Context) {
	var data [][]byte
	var result map[string]interface{}
	results := []map[string]interface{}{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("customers"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		// Set the cursor to the last key in the bucket
		cursor := b.Cursor()
		lastKey, _ := cursor.Last()

		// Iterate and add to data
		for k, v := cursor.Seek(lastKey); k != nil; k, v = cursor.Prev() {
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

func GetCustomer(c *gin.Context) {
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
