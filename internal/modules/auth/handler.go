package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/models"
	"github.com/haxxu/friend-flow-api/pkg/contextx"
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
			"ip":            ip,
			"ua":            ua,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.AuthService.Logout(req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   "Failed to revoke session: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Logged out successfully",
	})
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	uid, err := contextx.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   "invalid user ID",
		})
		return
	}

	if err := h.AuthService.LogoutAll(uid); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Logged out from all sessions",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	userID, err := h.AuthService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Issue new tokens
	accessToken, _ := h.AuthService.TokenService.GenerateAccessToken(userID)
	newRefreshToken, _ := h.AuthService.TokenService.GenerateRefreshToken(userID)

	// Replace the old session
	err = h.AuthService.AuthSessionService.Update(userID, req.RefreshToken, req.RefreshToken, req.UserAgent, req.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{"access_token": accessToken,
			"refresh_token": newRefreshToken,
		},
	})
}
