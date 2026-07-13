package handler

import (
	"net/http"

	"shopping-cart-service/internal/service"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) verifyUserAccess(c *gin.Context) bool {
	userId := c.Param("userId")
	requesterId := c.GetHeader("X-User-Id")
	if requesterId != "" && requesterId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return false
	}
	return true
}

func (h *CartHandler) GetCart(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")

	cart, err := h.svc.GetCart(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItemToCart(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")

	var req struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than zero"})
		return
	}

	cart, err := h.svc.AddItem(userId, req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")
	itemId := c.Param("itemId")

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity cannot be negative"})
		return
	}

	cart, err := h.svc.UpdateItem(userId, itemId, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) RemoveItemFromCart(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")
	itemId := c.Param("itemId")

	cart, err := h.svc.RemoveItem(userId, itemId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")

	err := h.svc.ClearCart(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

func (h *CartHandler) Checkout(c *gin.Context) {
	if !h.verifyUserAccess(c) {
		return
	}
	userId := c.Param("userId")

	order, err := h.svc.Checkout(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checkout successful! Order placed.",
		"order":   order,
	})
}
