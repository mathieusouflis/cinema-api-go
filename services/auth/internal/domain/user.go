package domain

import "context"

type User struct {
	ID       string
	Username string
	Email    string
	Password string
	GoogleID string
	Role     string
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateOAuthUserInput struct {
	Username string
	Email    string
	GoogleID string
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*User, error)
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	CreateWithOAuth(ctx context.Context, input CreateOAuthUserInput) (*User, error)
	UpdateGoogleID(ctx context.Context, userID string, googleID string) error
	Delete(ctx context.Context, id string) error
}
