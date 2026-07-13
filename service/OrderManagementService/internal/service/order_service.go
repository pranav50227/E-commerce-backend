package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"order-management-service/internal/model"
	"order-management-service/internal/repository"
)

type OrderService struct {
	repo         repository.OrderRepository
	productURL   string
	inventoryURL string
}

func NewOrderService(repo repository.OrderRepository, productURL string, inventoryURL string) *OrderService {
	return &OrderService{
		repo:         repo,
		productURL:   productURL,
		inventoryURL: inventoryURL,
	}
}

func (s *OrderService) GetUserOrders(userId string) ([]model.Order, error) {
	return s.repo.GetByUserID(userId)
}

func (s *OrderService) GetOrder(orderId string) (model.Order, error) {
	return s.repo.GetByID(orderId)
}

func (s *OrderService) Create(userId string, items []struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}) (model.Order, error) {
	var orderItems []model.OrderItem
	var totalPrice float64

	// Validate products & fetch prices
	for _, item := range items {
		resp, err := http.Get(fmt.Sprintf("%s/internal/products/%s", s.productURL, item.ProductID))
		if err != nil || resp.StatusCode != http.StatusOK {
			return model.Order{}, fmt.Errorf("product %s not found or invalid", item.ProductID)
		}
		defer resp.Body.Close()

		var product model.Product
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			return model.Order{}, errors.New("failed to decode product details")
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		totalPrice += product.Price * float64(item.Quantity)
	}

	// Deduct stock
	var deductReq struct {
		Items []struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}
	for _, item := range items {
		deductReq.Items = append(deductReq.Items, struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	bodyBytes, _ := json.Marshal(deductReq)
	resp, err := http.Post(s.inventoryURL+"/internal/inventory/deduct", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return model.Order{}, errors.New("failed to connect to InventoryService")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		respBody, _ := io.ReadAll(resp.Body)
		return model.Order{}, fmt.Errorf("failed to place order due to stock issues: %s", string(respBody))
	}

	if resp.StatusCode != http.StatusOK {
		return model.Order{}, errors.New("failed to deduct inventory")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	order := model.Order{
		ID:         fmt.Sprintf("ord-%d", r.Intn(1000000)),
		UserID:     userId,
		Items:      orderItems,
		TotalPrice: totalPrice,
		Status:     "Success",
		CreatedAt:  time.Now(),
	}

	err = s.repo.Save(order)
	return order, err
}

func (s *OrderService) UpdateStatus(orderId string, status string) (model.Order, error) {
	return s.repo.UpdateStatus(orderId, status)
}

func (s *OrderService) Cancel(userId string, orderId string) (model.Order, error) {
	order, err := s.repo.GetByID(orderId)
	if err != nil {
		return model.Order{}, err
	}

	if userId != "" && order.UserID != userId {
		return model.Order{}, errors.New("access denied")
	}

	if order.Status == "Cancelled" {
		return model.Order{}, errors.New("order is already cancelled")
	}

	// Restock
	for _, item := range order.Items {
		var restockReq = struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{ProductID: item.ProductID, Quantity: item.Quantity}
		bodyBytes, _ := json.Marshal(restockReq)
		_, _ = http.Post(s.inventoryURL+"/api/v1/inventory/restock", "application/json", bytes.NewBuffer(bodyBytes))
	}

	return s.repo.UpdateStatus(orderId, "Cancelled")
}
