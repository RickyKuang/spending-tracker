package store

import (
	"context"
	"time"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// CreateSession inserts a new session row for a user, storing only the hash of the token.
func (s *PostgresStore) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (models.Session, error) {
	// TODO: implement
	return models.Session{}, nil
}

// GetSessionByTokenHash looks up a session by SHA-256(token).
func (s *PostgresStore) GetSessionByTokenHash(ctx context.Context, tokenHash string) (models.Session, error) {
	// TODO: implement
	return models.Session{}, nil
}

// TouchSession updates last_used_at on a session.
func (s *PostgresStore) TouchSession(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}

// DeleteSession removes a session row, revoking it.
func (s *PostgresStore) DeleteSession(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}
