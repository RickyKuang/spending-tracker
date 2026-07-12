package apiclient

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// ListAccounts fetches all accounts across the user's linked institutions.
func (c *Client) ListAccounts(ctx context.Context) ([]models.Account, error) {
	// TODO: implement
	return nil, nil
}
