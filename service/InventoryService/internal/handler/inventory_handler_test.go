package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"inventory-service/internal/model"
	"inventory-service/internal/repository"
	"inventory-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestInventoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryInventoryRepository()
	svc := service.NewInventoryService(repo)
	h := NewInventoryHandler(svc)

	r := gin.Default()
	r.GET("/api/v1/inventory", h.GetInventory)
	r.GET("/api/v1/inventory/:productId", h.GetInventoryByProduct)
	r.PUT("/api/v1/inventory/:productId", h.UpdateInventory)
	r.POST("/api/v1/inventory/restock", h.RestockInventory)
	r.POST("/internal/inventory/deduct", h.DeductInventoryInternal)

	// Test GetInventory
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/inventory", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test GetInventoryByProduct
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/inventory/prod1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test UpdateInventory
	updateReq := map[string]int{"quantity": 150}
	body, _ := json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/inventory/prod1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test RestockInventory
	restockReq := map[string]interface{}{"productId": "prod1", "quantity": 10}
	body, _ = json.Marshal(restockReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/inventory/restock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test DeductInventoryInternal
	deductReq := model.DeductRequest{
		Items: []model.DeductItem{
			{ProductID: "prod1", Quantity: 5},
		},
	}
	body, _ = json.Marshal(deductReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/internal/inventory/deduct", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
