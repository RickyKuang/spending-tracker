package plaid

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/bank"
)

// CreateLinkToken starts a new Link session for a given user via /link/token/create.
func (c *Client) CreateLinkToken(ctx context.Context, clientUserID string) (bank.LinkTokenCreateResult, error) {
	// TODO: implement
	return bank.LinkTokenCreateResult{}, nil
}

// ExchangePublicToken exchanges a Link public token for a permanent access token via /item/public_token/exchange.
func (c *Client) ExchangePublicToken(ctx context.Context, publicToken string) (bank.ExchangeResult, error) {
	// TODO: implement
	return bank.ExchangeResult{}, nil
}
