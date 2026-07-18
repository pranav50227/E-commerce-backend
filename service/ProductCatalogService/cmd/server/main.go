package main

import (
	"log"
	"net/http"

	"product-catalog-service/api"
	"product-catalog-service/internal/handler"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"
	"github.com/gin-gonic/gin"

	"shared/utils"
)

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryProductRepository()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Product Catalog Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	port := utils.GetEnv("PORT_PRODUCT", "8081")
	log.Printf("Product Catalog Service starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Product Catalog Service: %v", err)
	}
}
