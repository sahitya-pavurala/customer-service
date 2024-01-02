package main

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Customer struct {
	CustomerID      string    `json:"customer_id"`
	Profile         string    `json:"profile"`
	PhysicalAddress string    `json:"physical_address"`
	Accounts        []Account `json:"accounts"`
}

type Account struct {
	AccountID string `json:"account_id"`
	Type      string `json:"account_type"`
}
