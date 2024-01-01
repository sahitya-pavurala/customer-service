package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"net/http"
  "time"
  "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("dummy_secret_key")

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context) {
	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		existingUser := b.Get([]byte(user.Username))
		if existingUser != nil {
			return fmt.Errorf("Username already exists")
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put([]byte(user.Username), jsonData)
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func LoginHandler(c *gin.Context) {
	var user User
	var existingUser User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return fmt.Errorf("Bucket not found")
		}

		storedUser := b.Get([]byte(user.Username))
		if storedUser == nil {
			return fmt.Errorf("User not found")
		}

		return json.Unmarshal(storedUser, &existingUser)
	})

	if err != nil || user.Password != existingUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		Subject:   user.Username,
		ExpiresAt: expirationTime.Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
