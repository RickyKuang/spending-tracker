package store

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// CreateUser inserts a new user with a pre-hashed password.
func (s *PostgresStore) CreateUser(ctx context.Context, email, passwordHash string) (models.User, error) {
	// TODO: implement
	return models.User{}, nil
}

// GetUserByEmail looks up a user by their (case-insensitive) email.
func (s *PostgresStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	// TODO: implement
	return models.User{}, nil
}

// GetUserByID looks up a user by primary key.
func (s *PostgresStore) GetUserByID(ctx context.Context, id string) (models.User, error) {
	// TODO: implement
	return models.User{}, nil
}
