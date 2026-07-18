package main

import (
	"log"
	"net/http"

	"auth-service/api"
	"auth-service/internal/handler"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"

	"shared/constants"
	"shared/utils"
)

func main() {
	r := gin.Default()

	// Initialize layers
	userSvcURL := utils.GetEnv("USER_SERVICE_URL", constants.UserServiceFallbackURL)
	svc := service.NewAuthService([]byte(constants.DefaultJWTSecret), userSvcURL)
	h := handler.NewAuthHandler(svc)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Auth Service is running"})
	})

	// Setup routes
	api.SetupRoutes(r, h)

	port := utils.GetEnv("PORT_AUTH", "8085")
	log.Printf("Auth Service starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Auth Service: %v", err)
	}
}
