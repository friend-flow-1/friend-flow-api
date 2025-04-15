package routes

import (
	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/auth"
)

func AuthRoutes(r *mux.Router, handler *auth.Handler) {
	r.HandleFunc("/auth/register", handler.Register).Methods("POST")
	r.HandleFunc("/auth/login", handler.Login).Methods("POST")
}
