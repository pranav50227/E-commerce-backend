package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/model"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up a mock UserManagementService server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/internal/users/" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(model.User{ID: "user-123"})
			return
		}
		if r.Method == "GET" && r.URL.Path == "/internal/users/username/john" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.User{
				ID:       "user-123",
				Username: "john",
				Password: "password123",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	jwtSecret := []byte("super-secret-key-12345")
	svc := service.NewAuthService(jwtSecret, server.URL)
	h := NewAuthHandler(svc)

	r := gin.Default()
	r.POST("/api/v1/auth/register", h.Register)
	r.POST("/api/v1/auth/login", h.Login)
	r.POST("/api/v1/auth/refresh", h.RefreshToken)

	// Test Register Handler
	regReq := map[string]string{
		"username": "john",
		"password": "password123",
		"name":     "John Doe",
		"email":    "john@example.com",
	}
	body, _ := json.Marshal(regReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected registration handler status 201, got %d", w.Code)
	}

	// Test Login Handler
	loginReq := map[string]string{
		"username": "john",
		"password": "password123",
	}
	body, _ = json.Marshal(loginReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected login handler status 200, got %d", w.Code)
	}

	var loginResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected non-empty token in login response")
	}

	// Test RefreshToken Handler
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/refresh", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected refresh handler status 200, got %d", w.Code)
	}
}
