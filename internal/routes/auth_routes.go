package routes

import (
	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
)

func AuthRoutes(r *mux.Router, handler *auth.AuthHandler) {
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
}
