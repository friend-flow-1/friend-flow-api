package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/haxxu/friend-flow-api/pkg/utils"
)

type AuthHandler struct {
	AuthService *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{AuthService: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterDTO
	log.Println(req)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	log.Println(req)

	// Validate the request payload
	if err := req.Validate(); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	user, err := h.AuthService.Register(req)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	// Validate the request payload
	if err := req.Validate(); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	token, _, err := h.AuthService.Login(req)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, token)
}
