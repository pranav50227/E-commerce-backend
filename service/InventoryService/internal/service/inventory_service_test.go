package service

import (
	"testing"
	"inventory-service/internal/model"
	"inventory-service/internal/repository"
)

func TestInventoryService(t *testing.T) {
	repo := repository.NewInMemoryInventoryRepository()
	svc := NewInventoryService(repo)

	// Test GetStockByProduct seeded product
	stock, err := svc.GetStockByProduct("prod1")
	if err != nil || stock != 100 {
		t.Fatalf("expected stock 100 for prod1, got: %d (error: %v)", stock, err)
	}

	// Test UpdateStock
	err = svc.UpdateStock("prod1", 120)
	if err != nil {
		t.Fatalf("failed to update stock: %v", err)
	}
	stock, _ = svc.GetStockByProduct("prod1")
	if stock != 120 {
		t.Errorf("expected stock 120, got %d", stock)
	}

	// Test RestockProduct
	newStock, err := svc.RestockProduct("prod1", 30)
	if err != nil || newStock != 150 {
		t.Errorf("expected stock 150, got %d (error: %v)", newStock, err)
	}

	// Test DeductStock success
	items := []model.DeductItem{
		{ProductID: "prod1", Quantity: 50},
		{ProductID: "prod2", Quantity: 10},
	}
	err = svc.DeductStock(items)
	if err != nil {
		t.Fatalf("expected successful deduction, got: %v", err)
	}

	stock, _ = svc.GetStockByProduct("prod1")
	if stock != 100 {
		t.Errorf("expected stock 100 after deduction, got %d", stock)
	}

	// Test DeductStock insufficient stock
	badItems := []model.DeductItem{
		{ProductID: "prod2", Quantity: 1000},
	}
	err = svc.DeductStock(badItems)
	if err == nil {
		t.Error("expected failure due to insufficient stock, got nil error")
	}
}
