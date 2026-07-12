package store

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// CreateInstitution inserts a newly linked institution with its encrypted access token.
func (s *PostgresStore) CreateInstitution(ctx context.Context, userID, plaidItemID string, accessTokenEnc []byte, plaidInstitutionID, name string) (models.Institution, error) {
	// TODO: implement
	return models.Institution{}, nil
}

// ListInstitutionsByUser lists all institutions linked by a user.
func (s *PostgresStore) ListInstitutionsByUser(ctx context.Context, userID string) ([]models.Institution, error) {
	// TODO: implement
	return nil, nil
}

// GetInstitution looks up a single institution by primary key.
func (s *PostgresStore) GetInstitution(ctx context.Context, id string) (models.Institution, error) {
	// TODO: implement
	return models.Institution{}, nil
}

// ClaimInstitutionForSync atomically transitions an institution from 'idle' to 'syncing'.
// found is false if the institution was not idle (a sync is already in flight).
func (s *PostgresStore) ClaimInstitutionForSync(ctx context.Context, id string) (institution models.Institution, found bool, err error) {
	// TODO: implement
	return models.Institution{}, false, nil
}

// UpdateInstitutionCursor persists the next sync cursor, committed alongside the batch it produced.
func (s *PostgresStore) UpdateInstitutionCursor(ctx context.Context, id, cursor string) error {
	// TODO: implement
	return nil
}

// FinishInstitutionSync marks an institution 'idle' after a successful sync.
func (s *PostgresStore) FinishInstitutionSync(ctx context.Context, id string) error {
	// TODO: implement
	return nil
}

// FailInstitutionSync marks an institution 'error' with a message, without failing the whole sync request.
func (s *PostgresStore) FailInstitutionSync(ctx context.Context, id, errMsg string) error {
	// TODO: implement
	return nil
}
