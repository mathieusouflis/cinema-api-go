package handler

import (
	registerUsecase "authService/internal/usecase/register"
	"encoding/json"
	"net/http"

	"filmserver/pkg/errors"
	"filmserver/pkg/render"
)

type RegisterHandler struct {
	usecase *registerUsecase.Usecase
}

func NewRegisterHandler(usecase *registerUsecase.Usecase) *RegisterHandler {
	return &RegisterHandler{usecase: usecase}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req registerUsecase.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.Render(w, errors.ErrBadRequest)
		return
	}

	_, err := h.usecase.Execute(r.Context(), req)

	if err != nil {
		errors.Render(w, err)
		return
	}

	render.Created(w, map[string]string{
		"message": "User created successfully",
	})
}
