package repository

import (
	"authService/db/orm"
	"authService/internal/domain"
	mappings "authService/internal/repository/postgres/utils/mappings"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresUserRepository struct {
	queries *orm.Queries
}

func NewPostgresUserRepository(queries *orm.Queries) *PostgresUserRepository {
	return &PostgresUserRepository{queries: queries}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return nil, err
	}
	user, err := r.queries.GetUserById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return mappings.MapUserToDomain(&user), nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return mappings.MapUserToDomain(&user), nil
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return mappings.MapUserToDomain(&user), nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
	user, err := r.queries.CreateUser(ctx, orm.CreateUserParams{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}
	return mappings.MapUserToDomain(&user), nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}
	return r.queries.DeleteUser(ctx, uuid)
}
