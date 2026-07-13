package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-management-service/internal/model"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	productServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.Product{ID: "prod1", Price: 10.00})
	}))
	defer productServer.Close()

	inventoryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer inventoryServer.Close()

	repo := repository.NewInMemoryOrderRepository()
	svc := service.NewOrderService(repo, productServer.URL, inventoryServer.URL)
	h := NewOrderHandler(svc)

	r := gin.Default()
	r.GET("/api/v1/orders", h.GetOrders)
	r.GET("/api/v1/orders/:orderId", h.GetOrderById)
	r.POST("/api/v1/orders", h.CreateOrder)
	r.PUT("/api/v1/orders/:orderId/status", h.UpdateOrderStatus)
	r.POST("/api/v1/orders/:orderId/cancel", h.CancelOrder)

	// Test CreateOrder
	createReq := map[string]interface{}{
		"userId": "user-1",
		"items": []map[string]interface{}{
			{"productId": "prod1", "quantity": 1},
		},
	}
	body, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/orders", bytes.NewBuffer(body))
	req.Header.Set("X-User-Id", "user-1")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var createdOrder model.Order
	_ = json.Unmarshal(w.Body.Bytes(), &createdOrder)

	// Test GetOrders
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/orders", nil)
	req.Header.Set("X-User-Id", "user-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test GetOrderById
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/orders/"+createdOrder.ID, nil)
	req.Header.Set("X-User-Id", "user-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test UpdateOrderStatus
	updateReq := map[string]string{"status": "Processing"}
	body, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/orders/"+createdOrder.ID+"/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test CancelOrder
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/orders/"+createdOrder.ID+"/cancel", nil)
	req.Header.Set("X-User-Id", "user-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
