package loginUsecase

import (
	"context"
	"time"

	"authService/internal/domain"

	"filmserver/pkg/errors"
	"filmserver/pkg/jwt"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type Input struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type Output struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Usecase struct {
	UserRepository  domain.UserRepository
	TokenRepository domain.TokenRepository
}

func New(userRepository domain.UserRepository, tokenRepository domain.TokenRepository) *Usecase {
	return &Usecase{
		UserRepository:  userRepository,
		TokenRepository: tokenRepository,
	}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	if input.Password == "" || input.UsernameOrEmail == "" {
		return Output{}, errors.ErrBadRequest
	}

	user, err := u.UserRepository.GetByEmail(ctx, input.UsernameOrEmail)
	if err != nil {
		user, err = u.UserRepository.GetByUsername(ctx, input.UsernameOrEmail)
		if err != nil {
			return Output{}, errors.ErrUnauth
		}
	}

	claims := jwt.JWTClaims{
		UserId: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}

	accessToken, err := jwt.SignToken(claims, AccessTokenTTL)
	if err != nil {
		return Output{}, err
	}

	refreshToken, err := jwt.SignToken(claims, RefreshTokenTTL)
	if err != nil {
		return Output{}, err
	}

	if err := u.TokenRepository.StoreRefreshToken(ctx, user.ID, refreshToken, RefreshTokenTTL); err != nil {
		return Output{}, err
	}

	return Output{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
