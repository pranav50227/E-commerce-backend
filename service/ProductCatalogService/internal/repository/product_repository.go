package repository

import (
	"errors"
	"strings"
	"sync"

	"product-catalog-service/internal/model"
)

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	GetByID(id string) (model.Product, error)
	Create(p model.Product) error
	Update(id string, update model.Product) (model.Product, error)
	Delete(id string) error
	GetByCategory(category string) ([]model.Product, error)
	Search(query string) ([]model.Product, error)
}

type InMemoryProductRepository struct {
	products map[string]model.Product
	mu       sync.RWMutex
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	repo := &InMemoryProductRepository{
		products: make(map[string]model.Product),
	}
	// Seed initial products
	repo.products["prod1"] = model.Product{ID: "prod1", Name: "Wireless Mouse", Description: "Ergonomic 2.4GHz wireless mouse", Price: 29.99, Category: "Electronics"}
	repo.products["prod2"] = model.Product{ID: "prod2", Name: "Mechanical Keyboard", Description: "RGB backlit mechanical keyboard", Price: 89.99, Category: "Electronics"}
	repo.products["prod3"] = model.Product{ID: "prod3", Name: "Water Bottle", Description: "Stainless steel vacuum insulated bottle", Price: 19.99, Category: "Lifestyle"}
	repo.products["prod4"] = model.Product{ID: "prod4", Name: "Running Shoes", Description: "Lightweight mesh running shoes", Price: 59.99, Category: "Footwear"}
	return repo
}

func (r *InMemoryProductRepository) GetAll() ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []model.Product
	for _, p := range r.products {
		list = append(list, p)
	}
	return list, nil
}

func (r *InMemoryProductRepository) GetByID(id string) (model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.products[id]
	if !exists {
		return model.Product{}, errors.New("product not found")
	}
	return p, nil
}

func (r *InMemoryProductRepository) Create(p model.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[p.ID]; exists {
		return errors.New("product ID already exists")
	}
	r.products[p.ID] = p
	return nil
}

func (r *InMemoryProductRepository) Update(id string, update model.Product) (model.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.products[id]
	if !exists {
		return model.Product{}, errors.New("product not found")
	}

	if update.Name != "" {
		p.Name = update.Name
	}
	if update.Description != "" {
		p.Description = update.Description
	}
	if update.Price > 0 {
		p.Price = update.Price
	}
	if update.Category != "" {
		p.Category = update.Category
	}

	r.products[id] = p
	return p, nil
}

func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(r.products, id)
	return nil
}

func (r *InMemoryProductRepository) GetByCategory(category string) ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []model.Product
	for _, p := range r.products {
		if strings.EqualFold(p.Category, category) {
			list = append(list, p)
		}
	}
	return list, nil
}

func (r *InMemoryProductRepository) Search(query string) ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []model.Product
	q := strings.ToLower(query)
	for _, p := range r.products {
		if strings.Contains(strings.ToLower(p.Name), q) || strings.Contains(strings.ToLower(p.Description), q) {
			list = append(list, p)
		}
	}
	return list, nil
}
