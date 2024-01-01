package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

var db *bolt.DB

func InitializeDB() error {
	var err error
	db, err = bolt.Open("customer.db", 0600, nil)
	if err != nil {
		return err
	}

	// Initialize customers buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("customers"))
		return err
	})

  // Initialize users buckets
  err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func SeedData() error {
	return db.Update(func(tx *bolt.Tx) error {
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
}
