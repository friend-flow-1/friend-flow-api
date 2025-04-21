package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	}

	user, err := h.UserService.GetUserById(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
		"error":   nil,
	})
}
