// Package sync orchestrates fetching transactions from a bank.Provider and persisting them to the store.
package sync

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/bank"
	"github.com/RickyKuang/spending-tracker/internal/crypto"
	"github.com/RickyKuang/spending-tracker/internal/store"
)

// InstitutionSyncResult summarizes one institution's sync outcome.
type InstitutionSyncResult struct {
	InstitutionID string
	Added         int
	Modified      int
	Removed       int
	Error         string
}

// Service orchestrates syncing all of a user's institutions. It backs both the HTTP
// sync handler and any future scheduled job — no HTTP concerns belong in this package.
type Service struct {
	store    *store.PostgresStore
	provider bank.Provider
	sealer   *crypto.Sealer
}

// New builds a sync Service.
func New(s *store.PostgresStore, provider bank.Provider, sealer *crypto.Sealer) *Service {
	return &Service{store: s, provider: provider, sealer: sealer}
}

// SyncUser syncs every institution belonging to userID, per the sync workflow contract:
// claim each idle institution, page through provider.SyncTransactions committing the
// cursor with each batch, refresh balances, then release the claim. A per-institution
// error is recorded and does not fail the rest of the request.
func (svc *Service) SyncUser(ctx context.Context, userID string) ([]InstitutionSyncResult, error) {
	// TODO: implement
	return nil, nil
}
