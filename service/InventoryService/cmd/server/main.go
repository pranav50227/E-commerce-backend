package main

import (
	"log"
	"net/http"

	"inventory-service/api"
	"inventory-service/internal/handler"
	"inventory-service/internal/repository"
	"inventory-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryInventoryRepository()
	svc := service.NewInventoryService(repo)
	h := handler.NewInventoryHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Inventory Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	log.Println("Inventory Service starting on port 8082...")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Failed to start Inventory Service: %v", err)
	}
}
