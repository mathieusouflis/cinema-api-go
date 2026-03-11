package repository

import (
	"authService/db/orm"
	"authService/internal/domain"
)

func MapUserToDomain(user *orm.User) *domain.User {
	return &domain.User{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}
}
