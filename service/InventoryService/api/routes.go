package api

import (
	"inventory-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers routes for the service
func SetupRoutes(r *gin.Engine, h *handler.InventoryHandler) {
	// Inventory routes
	public := r.Group("/api/v1/inventory")
	{
		public.GET("/", h.GetInventory)
		public.GET("/:productId", h.GetInventoryByProduct)
		public.PUT("/:productId", h.UpdateInventory)
		public.POST("/restock", h.RestockInventory)
	}

	// Internal APIs
	r.POST("/internal/inventory/deduct", h.DeductInventoryInternal)
}
