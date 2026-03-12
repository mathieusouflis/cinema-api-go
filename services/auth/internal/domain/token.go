package domain

import (
	"context"
	"time"

	"filmserver/pkg/pkg/jwt"
)

type Token struct {
	Token string `json:"token"`
}

type TokenPayload jwt.JWTClaims

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID string, token string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}
