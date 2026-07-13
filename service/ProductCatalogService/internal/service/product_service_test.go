package service

import (
	"testing"
	"product-catalog-service/internal/model"
	"product-catalog-service/internal/repository"
)

func TestProductService(t *testing.T) {
	repo := repository.NewInMemoryProductRepository()
	svc := NewProductService(repo)

	// Test GetAllProducts
	list, err := svc.GetAllProducts()
	if err != nil || len(list) != 4 {
		t.Fatalf("expected 4 products, got: %d", len(list))
	}

	// Test CreateProduct
	p := model.Product{
		ID:       "p-test",
		Name:     "Service Test Product",
		Price:    100.00,
		Category: "ServiceCategory",
	}
	err = svc.CreateProduct(p)
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	// Test GetProductByID
	res, err := svc.GetProductByID("p-test")
	if err != nil || res.Name != "Service Test Product" {
		t.Fatalf("failed to retrieve created product: %v", err)
	}

	// Test UpdateProduct
	updated, err := svc.UpdateProduct("p-test", model.Product{Name: "Updated Name"})
	if err != nil || updated.Name != "Updated Name" {
		t.Fatalf("failed to update product: %v", err)
	}

	// Test GetProductsByCategory
	catList, err := svc.GetProductsByCategory("ServiceCategory")
	if err != nil || len(catList) != 1 {
		t.Fatalf("expected 1 product by category, got: %d", len(catList))
	}

	// Test SearchProducts
	searchList, err := svc.SearchProducts("Updated")
	if err != nil || len(searchList) != 1 {
		t.Fatalf("expected 1 search match, got %d", len(searchList))
	}

	// Test DeleteProduct
	err = svc.DeleteProduct("p-test")
	if err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}
}
