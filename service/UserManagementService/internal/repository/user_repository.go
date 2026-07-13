package repository

import (
	"errors"
	"sync"

	"user-management-service/internal/model"
)

type UserRepository interface {
	Create(user model.User) error
	GetByID(id string) (model.User, error)
	GetByUsername(username string) (model.User, error)
	Update(id string, update model.User) (model.User, error)
	Delete(id string) error
	GetAll() ([]model.User, error)
}

type InMemoryUserRepository struct {
	users map[string]model.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		users: make(map[string]model.User),
	}
	// Seed initial test user
	repo.users["user1"] = model.User{
		ID:       "user1",
		Username: "testuser",
		Password: "password123",
		Name:     "Test User",
		Email:    "test@example.com",
	}
	return repo
}

func (r *InMemoryUserRepository) Create(user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Username == user.Username {
			return errors.New("username already exists")
		}
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByUsername(username string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (r *InMemoryUserRepository) Update(id string, update model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return model.User{}, errors.New("user not found")
	}

	if update.Name != "" {
		user.Name = update.Name
	}
	if update.Email != "" {
		user.Email = update.Email
	}
	r.users[id] = user
	return user, nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) GetAll() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []model.User
	for _, u := range r.users {
		list = append(list, u)
	}
	return list, nil
}
