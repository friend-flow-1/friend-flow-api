package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/middleware"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
)

func AuthRoutes(rg *gin.RouterGroup, handler *auth.AuthHandler) {
	rg.POST("/register", handler.Register)
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.RefreshToken)

	rg.POST("/logout", handler.Logout)

	auth := rg.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/logout-all", handler.LogoutAll)
	}
}
