package apiclient

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// ListTransactionsOptions holds the query parameters for GET /api/transactions.
type ListTransactionsOptions struct {
	AccountID string
	From      string
	To        string
	Limit     int
	Cursor    string
}

// ListTransactionsResult is the response from GET /api/transactions.
type ListTransactionsResult struct {
	Transactions []models.Transaction
	NextCursor   string
}

// ListTransactions fetches a page of transactions matching opts.
func (c *Client) ListTransactions(ctx context.Context, opts ListTransactionsOptions) (ListTransactionsResult, error) {
	// TODO: implement
	return ListTransactionsResult{}, nil
}
