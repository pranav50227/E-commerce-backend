package repository

import (
	"testing"
	"product-catalog-service/internal/model"
)

func TestInMemoryProductRepository(t *testing.T) {
	repo := NewInMemoryProductRepository()

	// Test GetAll
	list, err := repo.GetAll()
	if err != nil || len(list) != 4 {
		t.Fatalf("expected 4 seeded products, got: %d (error: %v)", len(list), err)
	}

	// Test GetByID
	p, err := repo.GetByID("prod1")
	if err != nil || p.Name != "Wireless Mouse" {
		t.Fatalf("expected to find prod1 (Wireless Mouse), got error: %v", err)
	}

	// Test Create
	newProd := model.Product{
		ID:          "prod-test",
		Name:        "Test Product",
		Description: "Testing repository",
		Price:       9.99,
		Category:    "TestCategory",
	}
	err = repo.Create(newProd)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	// Test GetByCategory
	list, err = repo.GetByCategory("TestCategory")
	if err != nil || len(list) != 1 {
		t.Errorf("expected 1 product under TestCategory, got %d (error: %v)", len(list), err)
	}

	// Test Search
	list, err = repo.Search("Testing")
	if err != nil || len(list) != 1 {
		t.Errorf("expected 1 product matching 'Testing', got %d", len(list))
	}

	// Test Update
	updateProd := model.Product{
		Price: 14.99,
	}
	updated, err := repo.Update("prod-test", updateProd)
	if err != nil || updated.Price != 14.99 {
		t.Errorf("failed to update price: %v", err)
	}

	// Test Delete
	err = repo.Delete("prod-test")
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}
}
