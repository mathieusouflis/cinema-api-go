package refreshUsecase

import (
	"context"

	"authService/internal/domain"
	loginUsecase "authService/internal/usecase/login"

	"filmserver/pkg/errors"
	"filmserver/pkg/pkg/jwt"
)

type Input struct {
	RefreshToken string
}

type Output struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Usecase struct {
	TokenRepository domain.TokenRepository
}

func New(tokenRepository domain.TokenRepository) *Usecase {
	return &Usecase{TokenRepository: tokenRepository}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	if input.RefreshToken == "" {
		return Output{}, errors.ErrUnauth
	}

	claims, err := jwt.ParseToken(input.RefreshToken)
	if err != nil {
		return Output{}, errors.ErrUnauth
	}

	stored, err := u.TokenRepository.GetRefreshToken(ctx, claims.UserId)
	if err != nil {
		return Output{}, err
	}
	if stored != input.RefreshToken {
		return Output{}, errors.ErrUnauth
	}

	tokenClaims := jwt.JWTClaims{
		UserId: claims.UserId,
		Email:  claims.Email,
		Role:   claims.Role,
	}

	accessToken, err := jwt.SignToken(tokenClaims, loginUsecase.AccessTokenTTL)
	if err != nil {
		return Output{}, err
	}

	newRefreshToken, err := jwt.SignToken(tokenClaims, loginUsecase.RefreshTokenTTL)
	if err != nil {
		return Output{}, err
	}

	if err := u.TokenRepository.StoreRefreshToken(ctx, claims.UserId, newRefreshToken, loginUsecase.RefreshTokenTTL); err != nil {
		return Output{}, err
	}

	return Output{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
