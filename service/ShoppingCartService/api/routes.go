package api

import (
	"shopping-cart-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers routes for the service
func SetupRoutes(r *gin.Engine, h *handler.CartHandler) {
	// Cart routes
	public := r.Group("/api/v1/cart")
	{
		public.GET("/:userId", h.GetCart)
		public.POST("/:userId/items", h.AddItemToCart)
		public.PUT("/:userId/items/:itemId", h.UpdateCartItem)
		public.DELETE("/:userId/items/:itemId", h.RemoveItemFromCart)
		public.DELETE("/:userId", h.ClearCart)
		public.POST("/:userId/checkout", h.Checkout)
	}
}
