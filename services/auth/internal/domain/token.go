package domain

import "filmserver/pkg/pkg/jwt"

type Token struct {
	Token string `json:"token"`
}

type TokenPayload jwt.JWTClaims
