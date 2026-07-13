package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth-service/internal/model"
)

func TestAuthService(t *testing.T) {
	jwtSecret := []byte("test-key")

	// Set up a mock UserManagementService server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.URL.Path == "/internal/users/" {
			var u model.User
			_ = json.NewDecoder(r.Body).Decode(&u)
			if u.Username == "duplicate" {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(u)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/internal/users/username/john" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(model.User{
				ID:       "user-1",
				Username: "john",
				Password: "password123",
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := NewAuthService(jwtSecret, server.URL)

	// Test GenerateToken
	token, err := svc.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("expected successful token generation, got error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	// Test RegisterUser success
	userId, err := svc.RegisterUser("john", "password123", "John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("expected registration to succeed, got: %v", err)
	}
	if userId == "" {
		t.Error("expected non-empty userId")
	}

	// Test RegisterUser duplicate
	_, err = svc.RegisterUser("duplicate", "password123", "Duplicate", "dup@example.com")
	if err == nil || err.Error() != "username already exists" {
		t.Errorf("expected duplicate username error, got: %v", err)
	}

	// Test LoginUser success
	tok, user, err := svc.LoginUser("john", "password123")
	if err != nil {
		t.Fatalf("expected successful login, got: %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected userId user-1, got %s", user.ID)
	}
	if tok == "" {
		t.Error("expected non-empty token")
	}

	// Test LoginUser wrong password
	_, _, err = svc.LoginUser("john", "wrongpass")
	if err == nil || err.Error() != "invalid username or password" {
		t.Errorf("expected invalid password error, got: %v", err)
	}
}
