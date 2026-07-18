package main

import (
	"log"
	"net/http"

	"user-management-service/api"
	"user-management-service/internal/handler"
	"user-management-service/internal/repository"
	"user-management-service/internal/service"
	"github.com/gin-gonic/gin"

	"shared/utils"
)

func main() {
	r := gin.Default()

	// Initialize layers
	repo := repository.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "User Management Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	port := utils.GetEnv("PORT_USER", "8080")
	log.Printf("User Management Service starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start User Management Service: %v", err)
	}
}
