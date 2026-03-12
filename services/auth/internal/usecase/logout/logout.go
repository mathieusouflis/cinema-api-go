package logoutUsecase

import (
	"context"

	"authService/internal/domain"

	"filmserver/pkg/errors"
	"filmserver/pkg/jwt"
)

type Input struct {
	RefreshToken string
}

type Usecase struct {
	TokenRepository domain.TokenRepository
}

func New(tokenRepository domain.TokenRepository) *Usecase {
	return &Usecase{TokenRepository: tokenRepository}
}

func (u *Usecase) Execute(ctx context.Context, input Input) error {
	if input.RefreshToken == "" {
		return errors.ErrUnauth
	}

	claims, err := jwt.ParseToken(input.RefreshToken)
	if err != nil {
		return errors.ErrUnauth
	}

	return u.TokenRepository.DeleteRefreshToken(ctx, claims.UserId)
}
