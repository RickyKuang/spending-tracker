package auth

import (
	"context"
	"time"

	"github.com/RickyKuang/spending-tracker/internal/models"
	"github.com/RickyKuang/spending-tracker/internal/store"
)

// SessionManager creates, validates, and revokes opaque bearer-token sessions.
type SessionManager struct {
	store *store.PostgresStore
	ttl   time.Duration
}

// NewSessionManager builds a SessionManager with a fixed session lifetime.
func NewSessionManager(s *store.PostgresStore, ttl time.Duration) *SessionManager {
	return &SessionManager{store: s, ttl: ttl}
}

// CreateSession generates a new random token, stores SHA-256(token), and returns the raw
// token to the caller exactly once.
func (m *SessionManager) CreateSession(ctx context.Context, userID string) (token string, session models.Session, err error) {
	// TODO: implement
	return "", models.Session{}, nil
}

// ValidateSession looks up a session by the raw bearer token and returns the owning user.
func (m *SessionManager) ValidateSession(ctx context.Context, token string) (models.User, error) {
	// TODO: implement
	return models.User{}, nil
}

// RevokeSession deletes a session row, invalidating its token immediately.
func (m *SessionManager) RevokeSession(ctx context.Context, token string) error {
	// TODO: implement
	return nil
}
