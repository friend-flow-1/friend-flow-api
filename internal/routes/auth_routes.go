package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
)

func AuthRoutes(rg *gin.RouterGroup, handler *auth.AuthHandler) {
	rg.POST("/register", handler.Register)
	rg.POST("/login", handler.Login)
}
