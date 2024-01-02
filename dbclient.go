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

	// Initialize users buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func SeedData() error {
	return db.Update(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))
		accountsBucket := tx.Bucket([]byte("accounts"))

		if customersBucket == nil || accountsBucket == nil {
			return fmt.Errorf("Buckets not found")
		}

		// Seed data
		seedData := map[string]Customer{
			"1": {
				CustomerID:      "1",
				Profile:         "Luffy",
				PhysicalAddress: "123 Main St",
				Accounts: []Account{
					{AccountID: "acc1", Type: "savings"},
					{AccountID: "acc2", Type: "checking"},
				},
			},
			"2": {
				CustomerID:      "2",
				Profile:         "Zoro",
				PhysicalAddress: "265 First St",
				Accounts: []Account{
					{AccountID: "acc3", Type: "savings"},
					{AccountID: "acc4", Type: "checking"},
				},
			},
		}

		for customerID, customer := range seedData {
			// Insert into "customers" bucket
			customerData, err := json.Marshal(customer)
			if err != nil {
				return err
			}
			if err := customersBucket.Put([]byte(customerID), customerData); err != nil {
				return err
			}

			// Insert associated accounts into "accounts" bucket
			for _, account := range customer.Accounts {
				accountData, err := json.Marshal(account)
				if err != nil {
					return err
				}

				if err := accountsBucket.Put([]byte(account.AccountID), accountData); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
