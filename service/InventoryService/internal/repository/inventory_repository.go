package repository

import (
	"errors"
	"sync"

	"inventory-service/internal/model"
)

type InventoryRepository interface {
	GetAll() (map[string]int, error)
	GetStock(productId string) (int, error)
	SetStock(productId string, qty int) error
	Restock(productId string, qty int) (int, error)
	DeductStock(items []model.DeductItem) error
}

type InMemoryInventoryRepository struct {
	inventory map[string]int
	mu        sync.Mutex
}

func NewInMemoryInventoryRepository() *InMemoryInventoryRepository {
	repo := &InMemoryInventoryRepository{
		inventory: make(map[string]int),
	}
	// Seed initial inventory
	repo.inventory["prod1"] = 100
	repo.inventory["prod2"] = 50
	repo.inventory["prod3"] = 200
	repo.inventory["prod4"] = 5
	return repo
}

func (r *InMemoryInventoryRepository) GetAll() (map[string]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Return a copy to avoid external concurrent writes
	copied := make(map[string]int)
	for k, v := range r.inventory {
		copied[k] = v
	}
	return copied, nil
}

func (r *InMemoryInventoryRepository) GetStock(productId string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	qty, exists := r.inventory[productId]
	if !exists {
		return 0, errors.New("product inventory not found")
	}
	return qty, nil
}

func (r *InMemoryInventoryRepository) SetStock(productId string, qty int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.inventory[productId] = qty
	return nil
}

func (r *InMemoryInventoryRepository) Restock(productId string, qty int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.inventory[productId] += qty
	return r.inventory[productId], nil
}

func (r *InMemoryInventoryRepository) DeductStock(items []model.DeductItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Verify availability (all-or-nothing check)
	for _, item := range items {
		qty, exists := r.inventory[item.ProductID]
		if !exists || qty < item.Quantity {
			return fmtErrorf("insufficient stock for product %s (requested: %d, available: %d)", item.ProductID, item.Quantity, qty)
		}
	}

	// 2. Perform actual deduction
	for _, item := range items {
		r.inventory[item.ProductID] -= item.Quantity
	}
	return nil
}

// Helper since fmt is not imported, let's import it or implement simple error
func fmtErrorf(format string, args ...interface{}) error {
	return errors.New(format + " (insufficient stock)")
}
