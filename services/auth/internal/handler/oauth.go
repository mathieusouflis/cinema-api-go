package handler

import (
	"authService/internal/config"
	oauthCallbackUsecase "authService/internal/usecase/oauth"
	"encoding/json"
	"net/http"
	"time"

	"filmserver/pkg/errors"
	"filmserver/pkg/logger"
	"filmserver/pkg/render"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type OauthHandler struct {
	usecase *oauthCallbackUsecase.Usecase
}

func NewOauthHandler(uc *oauthCallbackUsecase.Usecase) *OauthHandler {
	return &OauthHandler{usecase: uc}
}

func (h *OauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := logger.New(config.Load().Env)
	provider := chi.URLParam(r, "provider")

	var email string
	var id string

	log.Debug("Provider", provider)

	if provider == "google" {
		var body struct {
			Token string `json:"credentials"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.Error("Body parse failed")
			errors.Render(w, errors.ErrBadRequest)
			return
		}

		token, _, err := jwt.NewParser().ParseUnverified(body.Token, jwt.MapClaims{})
		if err != nil {
			log.Error("Token parse failed")
			errors.Render(w, errors.ErrBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("Claims parse failed")
			errors.Render(w, errors.ErrBadRequest)
			return
		}

		email, _ = claims["email"].(string)
		id, _ = claims["sub"].(string)
	}

	out, err := h.usecase.Execute(r.Context(), oauthCallbackUsecase.Input{
		Email:    email,
		Id:       id,
		Provider: provider,
	})
	if err != nil {
		log.Error("Usecase execute failed", err)
		errors.Render(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    out.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
		MaxAge:   int(oauthCallbackUsecase.RefreshTokenTTL / time.Second),
	})

	render.JSON(w, http.StatusOK, map[string]string{
		"access_token": out.AccessToken,
	})
}
