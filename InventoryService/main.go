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
		c.JSON(http.StatusOK, gin.H{"status": "Inventory Service is running"})
	})

	// Inventory routes
	inventory := r.Group("/api/v1/inventory")
	{
		inventory.GET("/", getInventory)
		inventory.GET("/:productId", getInventoryByProduct)
		inventory.PUT("/:productId", updateInventory)
		inventory.POST("/restock", restockInventory)
	}

	log.Println("Inventory Service starting on port 8082...")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to start Inventory Service: %v", err)
	}
}

func getInventory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get all inventory"})
}

func getInventoryByProduct(c *gin.Context) {
	productId := c.Param("productId")
	c.JSON(http.StatusOK, gin.H{"productId": productId, "message": "Get inventory for product"})
}

func updateInventory(c *gin.Context) {
	productId := c.Param("productId")
	c.JSON(http.StatusOK, gin.H{"productId": productId, "message": "Update inventory for product"})
}

func restockInventory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Restock inventory"})
}
