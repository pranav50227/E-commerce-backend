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

	"auth-service/internal/model"
	"shared/auth"
)

type AuthService struct {
	jwtSecret  []byte
	userSvcURL string
}

func NewAuthService(jwtSecret []byte, userSvcURL string) *AuthService {
	return &AuthService{
		jwtSecret:  jwtSecret,
		userSvcURL: userSvcURL,
	}
}

func (s *AuthService) GenerateToken(userId string) (string, error) {
	return auth.GenerateToken(userId, s.jwtSecret)
}

func (s *AuthService) RegisterUser(username, password, name, email string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	userId := fmt.Sprintf("user-%d", r.Intn(1000000))

	newUser := model.User{
		ID:       userId,
		Username: username,
		Password: password,
		Name:     name,
		Email:    email,
	}

	body, err := json.Marshal(newUser)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(s.userSvcURL+"/internal/users/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("failed to connect to UserManagementService")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return "", errors.New("username already exists")
	}

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("user registration failed internally: %s", string(respBody))
	}

	return userId, nil
}

func (s *AuthService) LoginUser(username, password string) (string, model.User, error) {
	url := fmt.Sprintf("%s/internal/users/username/%s", s.userSvcURL, username)
	resp, err := http.Get(url)
	if err != nil {
		return "", model.User{}, errors.New("failed to connect to UserManagementService")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", model.User{}, errors.New("invalid username or password")
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", model.User{}, errors.New("failed to decode user details")
	}

	if user.Password != password {
		return "", model.User{}, errors.New("invalid username or password")
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", model.User{}, errors.New("failed to generate token")
	}

	return token, user, nil
}
