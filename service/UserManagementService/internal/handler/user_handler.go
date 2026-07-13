package handler

import (
	"net/http"

	"user-management-service/internal/model"
	"user-management-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userId := c.Param("userId")
	requesterId := c.GetHeader("X-User-Id")

	if requesterId != "" && requesterId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	user, err := h.svc.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"email":    user.Email,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	requesterId := c.GetHeader("X-User-Id")

	if requesterId != "" && requesterId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.UpdateProfile(userId, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	requesterId := c.GetHeader("X-User-Id")

	if requesterId != "" && requesterId != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err := h.svc.DeleteAccount(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Internal endpoints
func (h *UserHandler) CreateUserInternal(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.Register(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByUsernameInternal(c *gin.Context) {
	username := c.Param("username")

	user, err := h.svc.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserByIDInternal(c *gin.Context) {
	userId := c.Param("userId")

	user, err := h.svc.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
