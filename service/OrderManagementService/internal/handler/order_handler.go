package handler

import (
	"net/http"

	"order-management-service/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userId := c.GetHeader("X-User-Id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header missing"})
		return
	}

	list, err := h.svc.GetUserOrders(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	orderId := c.Param("orderId")
	userId := c.GetHeader("X-User-Id")

	order, err := h.svc.GetOrder(orderId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if userId != "" && order.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userId := c.GetHeader("X-User-Id")
	var req struct {
		UserID string `json:"userId"`
		Items  []struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userId == "" {
		userId = req.UserID
	}

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order items cannot be empty"})
		return
	}

	order, err := h.svc.Create(userId, req.Items)
	if err != nil {
		if err.Error() != "" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderId := c.Param("orderId")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.UpdateStatus(orderId, req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	userId := c.GetHeader("X-User-Id")

	order, err := h.svc.Cancel(userId, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully", "order": order})
}
