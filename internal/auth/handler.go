package auth

import (
	"encoding/json"
	"net/http"

	"github.com/haxxu/friend-flow-api/pkg/utils"
)

type Handler struct {
	Service *AuthService
}

func NewHandler(service *AuthService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	user, err := h.Service.Register(req)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid payload")
		return
	}
	token, _, err := h.Service.Login(req)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, token)
}
