package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"user-management-service/internal/model"
	"user-management-service/internal/repository"
	"user-management-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryUserRepository()
	svc := service.NewUserService(repo)
	h := NewUserHandler(svc)

	r := gin.Default()
	r.GET("/api/v1/users/:userId", h.GetUserByID)
	r.PUT("/api/v1/users/:userId", h.UpdateUser)
	r.DELETE("/api/v1/users/:userId", h.DeleteUser)
	r.POST("/internal/users/", h.CreateUserInternal)

	// Test CreateUserInternal
	newUser := model.User{
		ID:       "user-h1",
		Username: "handleruser",
		Password: "password",
		Name:     "Handler User",
		Email:    "handler@example.com",
	}
	body, _ := json.Marshal(newUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/internal/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	// Test GetUserByID - Authorized (same user)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/users/user-h1", nil)
	req.Header.Set("X-User-Id", "user-h1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var respUser map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &respUser)
	if respUser["username"] != "handleruser" {
		t.Errorf("expected username handleruser, got %v", respUser["username"])
	}

	// Test GetUserByID - Forbidden (different user)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/users/user-h1", nil)
	req.Header.Set("X-User-Id", "other-user")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}

	// Test UpdateUser
	updateBody, _ := json.Marshal(model.User{Name: "New Name"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/users/user-h1", bytes.NewBuffer(updateBody))
	req.Header.Set("X-User-Id", "user-h1")
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Test DeleteUser
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/users/user-h1", nil)
	req.Header.Set("X-User-Id", "user-h1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
