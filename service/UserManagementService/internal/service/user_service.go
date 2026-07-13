package service

import (
	"user-management-service/internal/model"
	"user-management-service/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user model.User) error {
	return s.repo.Create(user)
}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByUsername(username string) (model.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) UpdateProfile(id string, update model.User) (model.User, error) {
	return s.repo.Update(id, update)
}

func (s *UserService) DeleteAccount(id string) error {
	return s.repo.Delete(id)
}
