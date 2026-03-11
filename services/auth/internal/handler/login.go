package handler

import (
	usecase "authService/internal/usecase/login"
	"encoding/json"
	"net/http"
	"time"

	"filmserver/pkg/errors"
	"filmserver/pkg/render"
)

type LoginHandler struct {
	usecase *usecase.Usecase
}

func NewLoginHandler(usecase *usecase.Usecase) *LoginHandler {
	return &LoginHandler{usecase: usecase}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req usecase.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.Render(w, errors.ErrBadRequest)
		return
	}

	output, err := h.usecase.Execute(r.Context(), req)
	if err != nil {
		errors.Render(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
		MaxAge:   int(usecase.RefreshTokenTTL / time.Second),
	})

	render.JSON(w, http.StatusOK, map[string]string{
		"access_token": output.AccessToken,
	})
}
