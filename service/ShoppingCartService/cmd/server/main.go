package main

import (
	"log"
	"net/http"
	"os"

	"shopping-cart-service/api"
	"shopping-cart-service/internal/handler"
	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"
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
	repo := repository.NewInMemoryCartRepository()
	productURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")
	inventoryURL := getEnv("INVENTORY_SERVICE_URL", "http://localhost:8082")
	orderURL := getEnv("ORDER_SERVICE_URL", "http://localhost:8083")
	svc := service.NewCartService(repo, productURL, inventoryURL, orderURL)
	h := handler.NewCartHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Shopping Cart Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	log.Println("Shopping Cart Service starting on port 8084...")
	if err := r.Run(":8084"); err != nil {
		log.Fatalf("Failed to start Shopping Cart Service: %v", err)
	}
}
