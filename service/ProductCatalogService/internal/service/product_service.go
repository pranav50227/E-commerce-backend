package service

import (
	"product-catalog-service/internal/model"
	"product-catalog-service/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id string) (model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) CreateProduct(p model.Product) error {
	return s.repo.Create(p)
}

func (s *ProductService) UpdateProduct(id string, update model.Product) (model.Product, error) {
	return s.repo.Update(id, update)
}

func (s *ProductService) DeleteProduct(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetProductsByCategory(category string) ([]model.Product, error) {
	return s.repo.GetByCategory(category)
}

func (s *ProductService) SearchProducts(query string) ([]model.Product, error) {
	return s.repo.Search(query)
}
