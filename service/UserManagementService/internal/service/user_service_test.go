package service

import (
	"testing"
	"user-management-service/internal/model"
	"user-management-service/internal/repository"
)

func TestUserService(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	svc := NewUserService(repo)

	// Test Registration
	newUser := model.User{
		ID:       "user-test",
		Username: "testreg",
		Password: "password",
		Name:     "Test Register",
		Email:    "reg@example.com",
	}

	err := svc.Register(newUser)
	if err != nil {
		t.Fatalf("expected registration to succeed, got %v", err)
	}

	// Test GetUserByID
	user, err := svc.GetUserByID("user-test")
	if err != nil {
		t.Fatalf("expected to find user, got error: %v", err)
	}
	if user.Username != "testreg" {
		t.Errorf("username mismatch, expected testreg, got %s", user.Username)
	}

	// Test GetUserByUsername
	user, err = svc.GetUserByUsername("testreg")
	if err != nil {
		t.Fatalf("expected to find user by username, got error: %v", err)
	}

	// Test UpdateProfile
	updatedUser := model.User{
		Name: "Updated Name",
	}
	updated, err := svc.UpdateProfile("user-test", updatedUser)
	if err != nil {
		t.Fatalf("expected successful update, got error: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("expected updated name to be Updated Name, got %s", updated.Name)
	}

	// Test DeleteAccount
	err = svc.DeleteAccount("user-test")
	if err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}

	_, err = svc.GetUserByID("user-test")
	if err == nil {
		t.Error("expected user to be deleted, got nil error")
	}
}
