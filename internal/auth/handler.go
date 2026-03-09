package auth

import (
	"net/http"

	"example.com/filmserver/pkg/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"token": "TODO"})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusCreated, map[string]string{"message": "registered"})
}
