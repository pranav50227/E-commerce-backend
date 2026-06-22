package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Shopping Cart Service is running"})
	})

	// Cart routes
	cart := r.Group("/api/v1/cart")
	{
		cart.GET("/:userId", getCart)
		cart.POST("/:userId/items", addItemToCart)
		cart.PUT("/:userId/items/:itemId", updateCartItem)
		cart.DELETE("/:userId/items/:itemId", removeItemFromCart)
		cart.DELETE("/:userId", clearCart)
		cart.POST("/:userId/checkout", checkout)
	}

	log.Println("Shopping Cart Service starting on port 8084...")
	if err := r.Run(":8084"); err != nil {
		log.Fatalf("Failed to start Shopping Cart Service: %v", err)
	}
}

func getCart(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "message": "Get cart for user"})
}

func addItemToCart(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusCreated, gin.H{"userId": userId, "message": "Item added to cart"})
}

func updateCartItem(c *gin.Context) {
	userId := c.Param("userId")
	itemId := c.Param("itemId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "itemId": itemId, "message": "Cart item updated"})
}

func removeItemFromCart(c *gin.Context) {
	userId := c.Param("userId")
	itemId := c.Param("itemId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "itemId": itemId, "message": "Item removed from cart"})
}

func clearCart(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "message": "Cart cleared"})
}

func checkout(c *gin.Context) {
	userId := c.Param("userId")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "message": "Checkout initiated"})
}
