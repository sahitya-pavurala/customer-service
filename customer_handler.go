package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("dummy_secret_key")

func GetAllCustomers(c *gin.Context) {
	var customers []Customer

	err := db.View(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))
		accountsBucket := tx.Bucket([]byte("accounts"))

		if customersBucket == nil || accountsBucket == nil {
			return fmt.Errorf("Buckets not found")
		}

		// Iterate through customers
		return customersBucket.ForEach(func(customerID, customerData []byte) error {
			var customer Customer
			if err := json.Unmarshal(customerData, &customer); err != nil {
				return err
			}

			// Fetch associated accounts
			for _, account := range customer.Accounts {
				accountData := accountsBucket.Get([]byte(account.AccountID))
				if accountData == nil {
					return fmt.Errorf("Account not found")
				}

				if err := json.Unmarshal(accountData, &account); err != nil {
					return err
				}
			}

			// Append customer to result
			customers = append(customers, customer)
			return nil
		})
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func GetCustomerById(c *gin.Context) {
	customerID := c.Param("customer_id")

	var customer Customer

	err := db.View(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))
		accountsBucket := tx.Bucket([]byte("accounts"))

		if customersBucket == nil || accountsBucket == nil {
			return fmt.Errorf("Buckets not found")
		}

		// Retrieve customer by customer_id
		customerData := customersBucket.Get([]byte(customerID))
		if customerData == nil {
			return fmt.Errorf("Customer not found")
		}

		// Unmarshal customer data
		if err := json.Unmarshal(customerData, &customer); err != nil {
			return err
		}

		// Fetch associated accounts
		for _, account := range customer.Accounts {
			accountData := accountsBucket.Get([]byte(account.AccountID))
			if accountData == nil {
				return fmt.Errorf("Account not found")
			}

			if err := json.Unmarshal(accountData, &account); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func GetCustomerByAccountId(c *gin.Context) {
	accountID := c.Param("account_id")

	var customer Customer

	err := db.View(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))
		accountsBucket := tx.Bucket([]byte("accounts"))

		if customersBucket == nil || accountsBucket == nil {
			return fmt.Errorf("Buckets not found")
		}

		// Iterate through customers to find the one with the specified account_id
		err := customersBucket.ForEach(func(customerID, customerData []byte) error {
			var currCustomer Customer
			if err := json.Unmarshal(customerData, &currCustomer); err != nil {
				return err
			}

			// Check if the account_id exists in the customer's accounts
			for _, account := range currCustomer.Accounts {
				if account.AccountID == accountID {
					customer = currCustomer

					// Fetch associated accounts
					for _, acc := range customer.Accounts {
						accountData := accountsBucket.Get([]byte(acc.AccountID))
						if accountData == nil {
							return fmt.Errorf("Account not found")
						}

						if err := json.Unmarshal(accountData, &acc); err != nil {
							return err
						}
					}

					return nil
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		if customer.CustomerID == "" {
			return fmt.Errorf("Customer not found for account_id: %s", accountID)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// AddCustomer handles the creation of a new customer
func AddCustomer(c *gin.Context) {
	var newCustomer Customer

	// Bind JSON request body to the Customer struct
	if err := c.BindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new customer to the database
	err := db.Update(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))

		if customersBucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		// Generate a unique customer ID (you may use a UUID library for better uniqueness)
		newCustomer.CustomerID = strconv.FormatInt(time.Now().UnixNano(), 10)

		// Marshal the customer data
		customerData, err := json.Marshal(newCustomer)
		if err != nil {
			return err
		}

		// Store the new customer in the database
		return customersBucket.Put([]byte(newCustomer.CustomerID), customerData)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newCustomer)
}

// AddAccount handles the creation of a new account for an existing customer
func AddAccount(c *gin.Context) {
	customerID := c.Param("id") // Retrieve customer ID from the URL parameter

	var newAccount Account

	// Bind JSON request body to the Account struct
	if err := c.BindJSON(&newAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new account to the customer's accounts in the database
	err := db.Update(func(tx *bolt.Tx) error {
		customersBucket := tx.Bucket([]byte("customers"))

		if customersBucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		// Retrieve the existing customer data
		existingCustomerData := customersBucket.Get([]byte(customerID))
		if existingCustomerData == nil {
			return fmt.Errorf("Customer not found")
		}

		// Unmarshal the existing customer data
		var existingCustomer Customer
		if err := json.Unmarshal(existingCustomerData, &existingCustomer); err != nil {
			return err
		}

		// Generate a unique account ID (you may use a UUID library for better uniqueness)
		newAccount.AccountID = strconv.FormatInt(time.Now().UnixNano(), 10)

		// Append the new account to the existing customer's accounts
		existingCustomer.Accounts = append(existingCustomer.Accounts, newAccount)

		// Marshal the updated customer data
		updatedCustomerData, err := json.Marshal(existingCustomer)
		if err != nil {
			return err
		}

		// Update the existing customer data in the database
		return customersBucket.Put([]byte(customerID), updatedCustomerData)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newAccount)
}
