package jwt

import (
	"time"

	"filmserver/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func SignToken(payload JWTClaims, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": payload.UserId,
		"email":   payload.Email,
		"role":    payload.Role,
		"exp":     time.Now().Add(expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyToken(token string) bool {
	_, err := ParseToken(token)
	return err == nil
}

func ParseToken(token string) (*JWTClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrUnauth
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return nil, errors.ErrUnauth
	}

	payload := &JWTClaims{
		UserId: claims["user_id"].(string),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
	}

	return payload, nil
}
