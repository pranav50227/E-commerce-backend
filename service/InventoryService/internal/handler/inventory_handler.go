package handler

import (
	"net/http"

	"inventory-service/internal/model"
	"inventory-service/internal/service"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
	inv, err := h.svc.GetInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inv)
}

func (h *InventoryHandler) GetInventoryByProduct(c *gin.Context) {
	productId := c.Param("productId")
	qty, err := h.svc.GetStockByProduct(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"productId": productId, "quantity": qty})
}

func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	productId := c.Param("productId")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.UpdateStock(productId, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"productId": productId, "quantity": req.Quantity, "message": "Inventory updated"})
}

func (h *InventoryHandler) RestockInventory(c *gin.Context) {
	var req struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ProductID == "" || req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID or quantity"})
		return
	}

	currentQty, err := h.svc.RestockProduct(req.ProductID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"productId": req.ProductID, "quantity": currentQty, "message": "Product restocked successfully"})
}

func (h *InventoryHandler) DeductInventoryInternal(c *gin.Context) {
	var req model.DeductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.DeductStock(req.Items)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory deducted successfully"})
}
