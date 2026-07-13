package main

import (
	"log"
	"net/http"
	"os"

	"order-management-service/api"
	"order-management-service/internal/handler"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"
	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryOrderRepository()
	productURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")
	inventoryURL := getEnv("INVENTORY_SERVICE_URL", "http://localhost:8082")
	svc := service.NewOrderService(repo, productURL, inventoryURL)
	h := handler.NewOrderHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Order Management Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	log.Println("Order Management Service starting on port 8083...")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Failed to start Order Management Service: %v", err)
	}
}
