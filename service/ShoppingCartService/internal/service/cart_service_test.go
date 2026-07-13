package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"shopping-cart-service/internal/repository"
)

func TestCartService(t *testing.T) {
	// Mock Product Service
	productServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/products/prod1" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer productServer.Close()

	// Mock Inventory Service
	inventoryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/inventory/prod1" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"quantity": 10}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer inventoryServer.Close()

	// Mock Order Service
	orderServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/api/v1/orders/" {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": "ord-123", "status": "Success"}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer orderServer.Close()

	repo := repository.NewInMemoryCartRepository()
	svc := NewCartService(repo, productServer.URL, inventoryServer.URL, orderServer.URL)

	// Test GetCart empty
	cart, err := svc.GetCart("user-1")
	if err != nil || len(cart.Items) != 0 {
		t.Fatalf("expected empty cart, got items: %d (error: %v)", len(cart.Items), err)
	}

	// Test AddItem success
	cart, err = svc.AddItem("user-1", "prod1", 2)
	if err != nil {
		t.Fatalf("expected successful item addition, got: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].ProductID != "prod1" || cart.Items[0].Quantity != 2 {
		t.Errorf("cart items mismatch, got: %v", cart.Items)
	}

	// Test AddItem insufficient stock
	_, err = svc.AddItem("user-1", "prod1", 20)
	if err == nil {
		t.Error("expected stock limit error, got nil")
	}

	// Test UpdateItem
	cart, err = svc.UpdateItem("user-1", "prod1", 5)
	if err != nil || cart.Items[0].Quantity != 5 {
		t.Errorf("failed to update item qty: %v", err)
	}

	// Test Checkout
	order, err := svc.Checkout("user-1")
	if err != nil {
		t.Fatalf("expected successful checkout, got: %v", err)
	}
	if order == nil {
		t.Fatal("expected order response, got nil")
	}

	// Cart should be empty after checkout
	cart, _ = svc.GetCart("user-1")
	if len(cart.Items) != 0 {
		t.Errorf("expected cart to be empty after checkout, got %d items", len(cart.Items))
	}
}
