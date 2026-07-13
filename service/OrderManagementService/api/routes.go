package api

import (
	"order-management-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers routes for the service
func SetupRoutes(r *gin.Engine, h *handler.OrderHandler) {
	// Order routes
	public := r.Group("/api/v1/orders")
	{
		public.GET("/", h.GetOrders)
		public.GET("/:orderId", h.GetOrderById)
		public.POST("/", h.CreateOrder)
		public.PUT("/:orderId/status", h.UpdateOrderStatus)
		public.DELETE("/:orderId", h.CancelOrder)
	}
}
