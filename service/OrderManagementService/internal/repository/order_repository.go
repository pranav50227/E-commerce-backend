package repository

import (
	"errors"
	"sync"

	"order-management-service/internal/model"
)

type OrderRepository interface {
	GetByUserID(userId string) ([]model.Order, error)
	GetByID(id string) (model.Order, error)
	Save(order model.Order) error
	UpdateStatus(id string, status string) (model.Order, error)
}

type InMemoryOrderRepository struct {
	orders map[string]model.Order
	mu     sync.RWMutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[string]model.Order),
	}
}

func (r *InMemoryOrderRepository) GetByUserID(userId string) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userOrders []model.Order
	for _, o := range r.orders {
		if o.UserID == userId {
			userOrders = append(userOrders, o)
		}
	}
	return userOrders, nil
}

func (r *InMemoryOrderRepository) GetByID(id string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

func (r *InMemoryOrderRepository) Save(order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) UpdateStatus(id string, status string) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[id]
	if !exists {
		return model.Order{}, errors.New("order not found")
	}

	order.Status = status
	r.orders[id] = order
	return order, nil
}
