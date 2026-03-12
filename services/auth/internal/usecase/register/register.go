package registerUsecase

import (
	"context"
	"net/mail"
	"regexp"
	"strings"

	"authService/internal/domain"

	"filmserver/pkg/errors"
)

type Input struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Output struct{}

type Usecase struct {
	UserRepository domain.UserRepository
}

func New(userRepository domain.UserRepository) *Usecase {
	return &Usecase{UserRepository: userRepository}
}

func (u *Usecase) Execute(ctx context.Context, input Input) (Output, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" {
		return Output{}, errors.ErrBadRequest
	}

	email, err := mail.ParseAddress(input.Email)
	if err != nil {
		return Output{}, errors.ErrBadRequest
	}

	username := strings.ToLower(input.Username)

	if len(username) <= 8 {
		return Output{}, errors.ErrBadRequest
	}

	usernameRegex := regexp.MustCompile(`^[a-z0-9_.]+$`)
	if !usernameRegex.MatchString(username) {
		return Output{}, errors.ErrBadRequest
	}

	// passwordRegex enforces a strong password policy:
	// - At least 8 characters long
	// - At least one uppercase letter (A-Z)
	// - At least one lowercase letter (a-z)
	// - At least one digit (0-9)
	// - At least one special character from: !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?
	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()\-_=+\[\]{};:'",.<>?/\\|]).{8,}$`)
	if !passwordRegex.MatchString(input.Password) {
		return Output{}, errors.ErrBadRequest
	}

	_, err = u.UserRepository.Create(ctx, domain.CreateUserInput{
		Username: username,
		Email:    email.Address,
		Password: input.Password,
	})

	return Output{}, err
}
