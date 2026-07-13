package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-management-service/internal/model"
	"order-management-service/internal/repository"
)

func TestOrderService(t *testing.T) {
	// Set up mock product service server
	productServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/products/prod1" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.Product{
				ID:    "prod1",
				Price: 20.00,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer productServer.Close()

	// Set up mock inventory service server
	inventoryServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/internal/inventory/deduct" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/api/v1/inventory/restock" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer inventoryServer.Close()

	repo := repository.NewInMemoryOrderRepository()
	svc := NewOrderService(repo, productServer.URL, inventoryServer.URL)

	// Test Create order success
	items := []struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}{
		{ProductID: "prod1", Quantity: 2},
	}
	order, err := svc.Create("user-1", items)
	if err != nil {
		t.Fatalf("expected successful order creation, got: %v", err)
	}
	if order.TotalPrice != 40.00 {
		t.Errorf("expected total price 40.00, got %f", order.TotalPrice)
	}

	// Test GetOrder
	retrieved, err := svc.GetOrder(order.ID)
	if err != nil || retrieved.ID != order.ID {
		t.Fatalf("failed to retrieve order: %v", err)
	}

	// Test GetUserOrders
	userOrders, err := svc.GetUserOrders("user-1")
	if err != nil || len(userOrders) != 1 {
		t.Errorf("expected 1 order for user-1, got %d (error: %v)", len(userOrders), err)
	}

	// Test UpdateStatus
	updated, err := svc.UpdateStatus(order.ID, "Shipped")
	if err != nil || updated.Status != "Shipped" {
		t.Errorf("failed to update order status: %v", err)
	}

	// Test Cancel order
	cancelled, err := svc.Cancel("user-1", order.ID)
	if err != nil || cancelled.Status != "Cancelled" {
		t.Errorf("failed to cancel order: %v", err)
	}
}
