package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
)

type Handlers struct {
	AuthHandler *auth.AuthHandler
	UserHandler *user.UserHandler
}

func SetupRoutes(r *gin.Engine, h *Handlers) http.Handler {
	api := r.Group("/api/v1")

	AuthRoutes(api.Group("/auth"), h.AuthHandler)

	UserRoutes(api.Group("/users"), h.UserHandler)

	MeRoutes(api.Group("/me"), h.UserHandler)

	return r
}
