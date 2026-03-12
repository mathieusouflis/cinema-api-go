package domain

import "context"

type User struct {
	ID       string
	Username string
	Email    string
	Password string
	Role     string
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Delete(ctx context.Context, id string) error
}
