package store

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// UpsertAccount inserts or updates an account keyed on plaid_account_id.
func (s *PostgresStore) UpsertAccount(ctx context.Context, account models.Account) (models.Account, error) {
	// TODO: implement
	return models.Account{}, nil
}

// ListAccountsByInstitution lists all accounts under one institution.
func (s *PostgresStore) ListAccountsByInstitution(ctx context.Context, institutionID string) ([]models.Account, error) {
	// TODO: implement
	return nil, nil
}

// ListAccountsByUser lists all accounts across every institution belonging to a user.
func (s *PostgresStore) ListAccountsByUser(ctx context.Context, userID string) ([]models.Account, error) {
	// TODO: implement
	return nil, nil
}
