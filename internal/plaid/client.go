// Package plaid implements bank.Provider using the official Plaid Go SDK.
package plaid

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/bank"
	plaidgo "github.com/plaid/plaid-go/v27/plaid"
)

// Client implements bank.Provider using the Plaid API.
type Client struct {
	api      *plaidgo.APIClient
	clientID string
	secret   string
	env      string
}

// New builds a Plaid Client from client credentials and an environment name (sandbox|development|production).
func New(clientID, secret, env string) *Client {
	// TODO: implement
	return nil
}

// SyncTransactions fetches the next page of transaction changes for an Item via /transactions/sync.
func (c *Client) SyncTransactions(ctx context.Context, accessToken string, cursor string) (bank.SyncResult, error) {
	// TODO: implement
	return bank.SyncResult{}, nil
}

// GetAccounts fetches current account and balance data for an Item via /accounts/get.
func (c *Client) GetAccounts(ctx context.Context, accessToken string) ([]bank.Account, error) {
	// TODO: implement
	return nil, nil
}
