package models

import "time"

// User represents a user profile
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// Product represents a catalog item
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID string     `json:"userId"`
	Items  []CartItem `json:"items"`
}

// OrderItem represents an item purchase record
type OrderItem struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Order represents a purchase transaction receipt
type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"totalPrice"`
	Status     string      `json:"status"` // "Pending", "Success", "Failed", "Cancelled"
	CreatedAt  time.Time   `json:"createdAt"`
}

// DeductItem represents an inventory deduction quantity
type DeductItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

// DeductRequest represents an inventory deduction request body
type DeductRequest struct {
	Items []DeductItem `json:"items"`
}
