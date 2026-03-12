package registerUsecase

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

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

var (
	usernameRegex = regexp.MustCompile(`^[a-z0-9_.]+$`)
	specialChars  = regexp.MustCompile(`[!@#$%^&*()\-_=+\[\]{};:'",.<>?/\\|]`)
)

// isStrongPassword checks: ≥8 chars, upper, lower, digit, special char.
// Uses plain checks instead of lookaheads (unsupported in Go RE2).
func isStrongPassword(p string) bool {
	if len(p) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range p {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit && specialChars.MatchString(p)
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

	if !usernameRegex.MatchString(username) {
		return Output{}, errors.ErrBadRequest
	}

	// Go's regexp (RE2) doesn't support lookaheads, so we check each
	// password policy condition separately.
	if !isStrongPassword(input.Password) {
		return Output{}, errors.ErrBadRequest
	}

	_, err = u.UserRepository.Create(ctx, domain.CreateUserInput{
		Username: username,
		Email:    email.Address,
		Password: input.Password,
	})

	return Output{}, err
}
