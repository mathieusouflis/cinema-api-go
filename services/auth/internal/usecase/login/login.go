package usecase

import (
	"context"
	"time"

	repository "authService/internal/repository/postgres"

	"filmserver/pkg/errors"
	"filmserver/pkg/pkg/jwt"
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
	UserRepository repository.PostgresUserRepository
}

func NewLoginUseCase(userRepository *repository.PostgresUserRepository) *Usecase {
	return &Usecase{UserRepository: *userRepository}
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

	return Output{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
