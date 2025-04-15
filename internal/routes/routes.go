package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/auth"
	"github.com/haxxu/friend-flow-api/internal/user"
)

func SetupRoutes(userHandler *user.Handler, authHandler *auth.Handler) http.Handler {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	AuthRoutes(api, authHandler)
	UserRoutes(api, userHandler)

	return router
}
