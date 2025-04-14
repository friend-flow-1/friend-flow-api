package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/user"
)

func SetupRoutes(userHandler *user.Handler) http.Handler {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	UserRoutes(api, userHandler)

	return router
}
