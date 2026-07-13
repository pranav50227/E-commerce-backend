package api

import (
	"user-management-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers routes for the service
func SetupRoutes(r *gin.Engine, h *handler.UserHandler) {
	// Public routes
	publicUsers := r.Group("/api/v1/users")
	{
		publicUsers.GET("/:userId", h.GetUserByID)
		publicUsers.PUT("/:userId", h.UpdateUser)
		publicUsers.DELETE("/:userId", h.DeleteUser)
	}

	// Internal routes
	internal := r.Group("/internal/users")
	{
		internal.POST("/", h.CreateUserInternal)
		internal.GET("/username/:username", h.GetUserByUsernameInternal)
		internal.GET("/id/:userId", h.GetUserByIDInternal)
	}
}
