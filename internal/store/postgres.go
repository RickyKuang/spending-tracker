// Package store contains all Postgres query methods, grouped by aggregate.
package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore wraps a pgx connection pool and exposes query methods per aggregate.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPool opens a pgx connection pool against databaseURL.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}

// New builds a PostgresStore around an already-open pool.
func New(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

// ResetStuckSyncs resets any institution left in 'syncing' status back to 'idle'.
// Call once at startup: a process killed mid-sync leaves rows stuck, and since the cursor
// only commits alongside its batch, there is no partial state this could overwrite.
func (s *PostgresStore) ResetStuckSyncs(ctx context.Context) error {
	// TODO: implement
	return nil
}
