package service

import (
	"inventory-service/internal/model"
	"inventory-service/internal/repository"
)

type InventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetInventory() (map[string]int, error) {
	return s.repo.GetAll()
}

func (s *InventoryService) GetStockByProduct(productId string) (int, error) {
	return s.repo.GetStock(productId)
}

func (s *InventoryService) UpdateStock(productId string, qty int) error {
	return s.repo.SetStock(productId, qty)
}

func (s *InventoryService) RestockProduct(productId string, qty int) (int, error) {
	return s.repo.Restock(productId, qty)
}

func (s *InventoryService) DeductStock(items []model.DeductItem) error {
	return s.repo.DeductStock(items)
}
