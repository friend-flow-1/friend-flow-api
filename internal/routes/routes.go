package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
)

type Handlers struct {
	AuthHandler *auth.AuthHandler
	UserHandler *user.UserHandler
}

func SetupRoutes(r *mux.Router, h *Handlers) http.Handler {
	api := r.PathPrefix("/api/v1").Subrouter()

	AuthRoutes(api, h.AuthHandler)
	UserRoutes(api, h.UserHandler)

	return r
}
