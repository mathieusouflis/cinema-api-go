package handler

import (
	refreshUsecase "authService/internal/usecase/refresh"
	"net/http"
	"time"

	loginUsecase "authService/internal/usecase/login"

	"filmserver/pkg/errors"
	"filmserver/pkg/render"
)

type RefreshHandler struct {
	usecase *refreshUsecase.Usecase
}

func NewRefreshHandler(uc *refreshUsecase.Usecase) *RefreshHandler {
	return &RefreshHandler{usecase: uc}
}

func (h *RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		errors.Render(w, errors.ErrUnauth)
		return
	}

	output, err := h.usecase.Execute(r.Context(), refreshUsecase.Input{
		RefreshToken: cookie.Value,
	})
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
		MaxAge:   int(loginUsecase.RefreshTokenTTL / time.Second),
	})

	render.JSON(w, http.StatusOK, map[string]string{
		"access_token": output.AccessToken,
	})
}
