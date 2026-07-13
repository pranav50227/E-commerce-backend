package main

import (
	"log"
	"net/http"

	"user-management-service/api"
	"user-management-service/internal/handler"
	"user-management-service/internal/repository"
	"user-management-service/internal/service"
	"github.com/gin-gonic/gin"
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

	log.Println("User Management Service starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start User Management Service: %v", err)
	}
}
