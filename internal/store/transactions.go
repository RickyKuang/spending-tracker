package store

import (
	"context"
	"time"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// TransactionFilter narrows ListTransactions by account, date range, and pagination cursor.
type TransactionFilter struct {
	AccountID string
	From      *time.Time
	To        *time.Time
	Limit     int
	Cursor    string
}

// UpsertTransaction inserts or updates a transaction keyed on plaid_transaction_id.
// added and modified are treated uniformly, per the sync workflow's idempotency contract.
func (s *PostgresStore) UpsertTransaction(ctx context.Context, txn models.Transaction) (models.Transaction, error) {
	// TODO: implement
	return models.Transaction{}, nil
}

// MarkTransactionsRemoved soft-deletes transactions by plaid_transaction_id.
func (s *PostgresStore) MarkTransactionsRemoved(ctx context.Context, plaidTransactionIDs []string) error {
	// TODO: implement
	return nil
}

// ListTransactions lists live transactions for a user, filtered and paginated.
func (s *PostgresStore) ListTransactions(ctx context.Context, userID string, filter TransactionFilter) (transactions []models.Transaction, nextCursor string, err error) {
	// TODO: implement
	return nil, "", nil
}
