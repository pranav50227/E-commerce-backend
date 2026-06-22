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
		c.JSON(http.StatusOK, gin.H{"status": "Order Management Service is running"})
	})

	// Order routes
	orders := r.Group("/api/v1/orders")
	{
		orders.GET("/", getOrders)
		orders.GET("/:orderId", getOrderById)
		orders.POST("/", createOrder)
		orders.PUT("/:orderId/status", updateOrderStatus)
		orders.DELETE("/:orderId", cancelOrder)
	}

	log.Println("Order Management Service starting on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Failed to start Order Management Service: %v", err)
	}
}

func getOrders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all orders"})
}

func getOrderById(c *gin.Context) {
	orderId := c.Param("orderId")
	c.JSON(http.StatusOK, gin.H{"orderId": orderId, "message": "Get order by ID"})
}

func createOrder(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

func updateOrderStatus(c *gin.Context) {
	orderId := c.Param("orderId")
	c.JSON(http.StatusOK, gin.H{"orderId": orderId, "message": "Order status updated"})
}

func cancelOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	c.JSON(http.StatusOK, gin.H{"orderId": orderId, "message": "Order cancelled"})
}
