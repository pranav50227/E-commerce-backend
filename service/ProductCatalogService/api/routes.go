package api

import (
	"product-catalog-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers routes for the service
func SetupRoutes(r *gin.Engine, h *handler.ProductHandler) {
	// Product routes
	public := r.Group("/api/v1/products")
	{
		public.GET("/", h.GetProducts)
		public.GET("/search", h.SearchProducts)
		public.GET("/category/:category", h.GetProductsByCategory)
		public.GET("/:productId", h.GetProductByID)
		public.POST("/", h.CreateProduct)
		public.PUT("/:productId", h.UpdateProduct)
		public.DELETE("/:productId", h.DeleteProduct)
	}

	// Internal validation endpoint
	r.GET("/internal/products/:productId", h.GetProductByIdInternal)
}
