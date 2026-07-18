package main

import (
	"log"
	"net/http"

	"order-management-service/api"
	"order-management-service/internal/handler"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"
	"github.com/gin-gonic/gin"

	"shared/utils"
)

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryOrderRepository()
	productURL := utils.GetEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")
	inventoryURL := utils.GetEnv("INVENTORY_SERVICE_URL", "http://localhost:8082")
	svc := service.NewOrderService(repo, productURL, inventoryURL)
	h := handler.NewOrderHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Order Management Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	port := utils.GetEnv("PORT_ORDER", "8083")
	log.Printf("Order Management Service starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Order Management Service: %v", err)
	}
}
