package routes

import (
	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/user"
)

func UserRoutes(r *mux.Router, handler *user.Handler) {
	r.HandleFunc("/users/register", handler.Register).Methods("POST")
	r.HandleFunc("/users/login", handler.Login).Methods("POST")
}
