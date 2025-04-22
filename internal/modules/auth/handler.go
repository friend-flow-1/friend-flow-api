package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{AuthService: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.AuthService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data:    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Extract user agent and IP
	ua := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	// Call the login logic
	accessToken, refreshToken, user, err := h.AuthService.Login(req, ua, ip)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Return both tokens and user info
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user":          user,
		},
	})
}
