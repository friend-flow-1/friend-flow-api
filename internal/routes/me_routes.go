package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/middleware"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
)

func MeRoutes(rg *gin.RouterGroup, handler *user.UserHandler) {
	protected := rg.Group("/", middleware.AuthMiddleware())

	protected.GET("info", handler.GetUserInfo)
}
