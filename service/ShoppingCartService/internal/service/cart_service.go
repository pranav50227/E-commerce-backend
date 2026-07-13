package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"shopping-cart-service/internal/model"
	"shopping-cart-service/internal/repository"
)

type CartService struct {
	repo         repository.CartRepository
	productURL   string
	inventoryURL string
	orderURL     string
}

func NewCartService(repo repository.CartRepository, productURL, inventoryURL, orderURL string) *CartService {
	return &CartService{
		repo:         repo,
		productURL:   productURL,
		inventoryURL: inventoryURL,
		orderURL:     orderURL,
	}
}

func (s *CartService) GetCart(userId string) (model.Cart, error) {
	return s.repo.GetByUserID(userId)
}

func (s *CartService) AddItem(userId string, productId string, qty int) (model.Cart, error) {
	// 1. Verify product exists
	pResp, err := http.Get(fmt.Sprintf("%s/internal/products/%s", s.productURL, productId))
	if err != nil || pResp.StatusCode != http.StatusOK {
		return model.Cart{}, errors.New("product not found or invalid")
	}
	pResp.Body.Close()

	// 2. Verify stock
	iResp, err := http.Get(fmt.Sprintf("%s/api/v1/inventory/%s", s.inventoryURL, productId))
	if err != nil || iResp.StatusCode != http.StatusOK {
		return model.Cart{}, errors.New("failed to verify inventory")
	}
	defer iResp.Body.Close()

	var stock struct {
		Quantity int `json:"quantity"`
	}
	_ = json.NewDecoder(iResp.Body).Decode(&stock)

	if stock.Quantity < qty {
		return model.Cart{}, fmt.Errorf("insufficient stock (available: %d)", stock.Quantity)
	}

	cart, err := s.repo.GetByUserID(userId)
	if err != nil {
		return model.Cart{}, err
	}

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productId {
			if stock.Quantity < item.Quantity+qty {
				return model.Cart{}, fmt.Errorf("insufficient stock for combined quantity (available: %d)", stock.Quantity)
			}
			cart.Items[i].Quantity += qty
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, model.CartItem{ProductID: productId, Quantity: qty})
	}

	err = s.repo.Save(cart)
	return cart, err
}

func (s *CartService) UpdateItem(userId string, productId string, qty int) (model.Cart, error) {
	cart, err := s.repo.GetByUserID(userId)
	if err != nil {
		return model.Cart{}, err
	}

	foundIdx := -1
	for i, item := range cart.Items {
		if item.ProductID == productId {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return model.Cart{}, errors.New("item not found in cart")
	}

	if qty == 0 {
		cart.Items = append(cart.Items[:foundIdx], cart.Items[foundIdx+1:]...)
	} else {
		// Check inventory
		iResp, err := http.Get(fmt.Sprintf("%s/api/v1/inventory/%s", s.inventoryURL, productId))
		if err == nil && iResp.StatusCode == http.StatusOK {
			var stock struct {
				Quantity int `json:"quantity"`
			}
			_ = json.NewDecoder(iResp.Body).Decode(&stock)
			iResp.Body.Close()
			if stock.Quantity < qty {
				return model.Cart{}, fmt.Errorf("insufficient stock (available: %d)", stock.Quantity)
			}
		}
		cart.Items[foundIdx].Quantity = qty
	}

	err = s.repo.Save(cart)
	return cart, err
}

func (s *CartService) RemoveItem(userId string, productId string) (model.Cart, error) {
	cart, err := s.repo.GetByUserID(userId)
	if err != nil {
		return model.Cart{}, err
	}

	foundIdx := -1
	for i, item := range cart.Items {
		if item.ProductID == productId {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return model.Cart{}, errors.New("item not found in cart")
	}

	cart.Items = append(cart.Items[:foundIdx], cart.Items[foundIdx+1:]...)
	err = s.repo.Save(cart)
	return cart, err
}

func (s *CartService) ClearCart(userId string) error {
	return s.repo.Delete(userId)
}

func (s *CartService) Checkout(userId string) (interface{}, error) {
	cart, err := s.repo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Place order
	var orderReq = struct {
		UserID string `json:"userId"`
		Items  []struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
	}{UserID: userId}

	for _, item := range cart.Items {
		orderReq.Items = append(orderReq.Items, struct {
			ProductID string `json:"productId"`
			Quantity  int    `json:"quantity"`
		}{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	bodyBytes, _ := json.Marshal(orderReq)
	req, err := http.NewRequest("POST", s.orderURL+"/api/v1/orders/", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to connect to OrderManagementService during checkout")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("checkout failed during order placement: %s", string(respBody))
	}

	// Success: clear cart
	_ = s.repo.Delete(userId)

	var orderResponse interface{}
	_ = json.NewDecoder(resp.Body).Decode(&orderResponse)
	return orderResponse, nil
}
