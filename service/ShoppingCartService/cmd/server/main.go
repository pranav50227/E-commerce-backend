package main

import (
	"log"
	"net/http"

	"shopping-cart-service/api"
	"shopping-cart-service/internal/handler"
	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"
	"github.com/gin-gonic/gin"

	"shared/utils"
)

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryCartRepository()
	productURL := utils.GetEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")
	inventoryURL := utils.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8082")
	orderURL := utils.GetEnv("ORDER_SERVICE_URL", "http://localhost:8083")
	svc := service.NewCartService(repo, productURL, inventoryURL, orderURL)
	h := handler.NewCartHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Shopping Cart Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	port := utils.GetEnv("PORT_CART", "8084")
	log.Printf("Shopping Cart Service starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Shopping Cart Service: %v", err)
	}
}
