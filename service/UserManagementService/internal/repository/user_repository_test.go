package repository

import (
	"testing"
	"user-management-service/internal/model"
)

func TestInMemoryUserRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()

	// Test GetByID for default seeded user
	user, err := repo.GetByID("user1")
	if err != nil {
		t.Fatalf("expected to find seeded user1, got error: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", user.Username)
	}

	// Test Create User
	newUser := model.User{
		ID:       "user2",
		Username: "newuser",
		Password: "password",
		Name:     "New User",
		Email:    "new@example.com",
	}
	err = repo.Create(newUser)
	if err != nil {
		t.Fatalf("expected to create user, got error: %v", err)
	}

	// Test Create duplicate username
	err = repo.Create(newUser)
	if err == nil {
		t.Error("expected error when creating duplicate username, got nil")
	}

	// Test GetByUsername
	user, err = repo.GetByUsername("newuser")
	if err != nil {
		t.Fatalf("expected to find newuser, got error: %v", err)
	}
	if user.ID != "user2" {
		t.Errorf("expected ID 'user2', got '%s'", user.ID)
	}

	// Test Update
	updateUser := model.User{
		Name:  "Updated User",
		Email: "updated@example.com",
	}
	updated, err := repo.Update("user2", updateUser)
	if err != nil {
		t.Fatalf("expected successful update, got error: %v", err)
	}
	if updated.Name != "Updated User" || updated.Email != "updated@example.com" {
		t.Errorf("update fields mismatch, got name: %s, email: %s", updated.Name, updated.Email)
	}

	// Test Delete
	err = repo.Delete("user2")
	if err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}

	_, err = repo.GetByID("user2")
	if err == nil {
		t.Error("expected user2 to be deleted and return error, got nil")
	}
}
