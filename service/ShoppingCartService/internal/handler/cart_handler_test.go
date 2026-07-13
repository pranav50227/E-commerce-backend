package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shopping-cart-service/internal/repository"
	"shopping-cart-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestCartHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	productServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer productServer.Close()

	inventoryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"quantity": 10}`))
	}))
	defer inventoryServer.Close()

	orderServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id": "ord-123"}`))
	}))
	defer orderServer.Close()

	repo := repository.NewInMemoryCartRepository()
	svc := service.NewCartService(repo, productServer.URL, inventoryServer.URL, orderServer.URL)
	h := NewCartHandler(svc)

	r := gin.Default()
	r.GET("/api/v1/cart/:userId", h.GetCart)
	r.POST("/api/v1/cart/:userId/items", h.AddItemToCart)
	r.PUT("/api/v1/cart/:userId/items/:itemId", h.UpdateCartItem)
	r.DELETE("/api/v1/cart/:userId/items/:itemId", h.RemoveItemFromCart)
	r.POST("/api/v1/cart/:userId/checkout", h.Checkout)

	// Test GetCart
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/cart/user-1", nil)
	req.Header.Set("X-User-Id", "user-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test AddItemToCart
	addReq := map[string]interface{}{"productId": "prod1", "quantity": 2}
	body, _ := json.Marshal(addReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/cart/user-1/items", bytes.NewBuffer(body))
	req.Header.Set("X-User-Id", "user-1")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Test UpdateCartItem
	updateReq := map[string]interface{}{"quantity": 3}
	body, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/cart/user-1/items/prod1", bytes.NewBuffer(body))
	req.Header.Set("X-User-Id", "user-1")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test Checkout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/cart/user-1/checkout", nil)
	req.Header.Set("X-User-Id", "user-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
