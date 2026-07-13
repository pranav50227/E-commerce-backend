package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"product-catalog-service/internal/model"
	"product-catalog-service/internal/repository"
	"product-catalog-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestProductHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryProductRepository()
	svc := service.NewProductService(repo)
	h := NewProductHandler(svc)

	r := gin.Default()
	r.GET("/api/v1/products", h.GetProducts)
	r.GET("/api/v1/products/:productId", h.GetProductByID)
	r.POST("/api/v1/products", h.CreateProduct)
	r.PUT("/api/v1/products/:productId", h.UpdateProduct)
	r.DELETE("/api/v1/products/:productId", h.DeleteProduct)
	r.GET("/api/v1/products/category/:category", h.GetProductsByCategory)
	r.GET("/api/v1/products/search", h.SearchProducts)

	// Test GetProducts
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test CreateProduct
	newP := model.Product{
		ID:          "p-handler",
		Name:        "Handler Product",
		Description: "Descr",
		Price:       15.00,
		Category:    "Tech",
	}
	body, _ := json.Marshal(newP)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Test GetProductByID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/products/p-handler", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test GetProductsByCategory
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/products/category/Tech", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test SearchProducts
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/products/search?q=Handler", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test UpdateProduct
	updateBody, _ := json.Marshal(model.Product{Price: 20.00})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/products/p-handler", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test DeleteProduct
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/products/p-handler", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
