package repository

import (
	"sync"

	"shopping-cart-service/internal/model"
)

type CartRepository interface {
	GetByUserID(userId string) (model.Cart, error)
	Save(cart model.Cart) error
	Delete(userId string) error
}

type InMemoryCartRepository struct {
	carts map[string]model.Cart
	mu    sync.RWMutex
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	return &InMemoryCartRepository{
		carts: make(map[string]model.Cart),
	}
}

func (r *InMemoryCartRepository) GetByUserID(userId string) (model.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart, exists := r.carts[userId]
	if !exists {
		return model.Cart{UserID: userId, Items: []model.CartItem{}}, nil
	}
	return cart, nil
}

func (r *InMemoryCartRepository) Save(cart model.Cart) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.carts[cart.UserID] = cart
	return nil
}

func (r *InMemoryCartRepository) Delete(userId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.carts, userId)
	return nil
}
