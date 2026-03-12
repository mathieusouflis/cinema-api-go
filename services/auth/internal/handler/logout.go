package handler

import (
	logoutUsecase "authService/internal/usecase/logout"
	"net/http"

	"filmserver/pkg/errors"
	"filmserver/pkg/render"
)

type LogoutHandler struct {
	usecase *logoutUsecase.Usecase
}

func NewLogoutHandler(uc *logoutUsecase.Usecase) *LogoutHandler {
	return &LogoutHandler{usecase: uc}
}

func (h *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		errors.Render(w, errors.ErrUnauth)
		return
	}

	if err := h.usecase.Execute(r.Context(), logoutUsecase.Input{
		RefreshToken: cookie.Value,
	}); err != nil {
		errors.Render(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
		MaxAge:   -1,
	})

	render.NoContent(w)
}
