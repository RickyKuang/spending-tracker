package apiclient

import (
	"context"

	"github.com/RickyKuang/spending-tracker/internal/models"
)

// LinkTokenCreateResult is the response from POST /api/plaid/link/token.
type LinkTokenCreateResult struct {
	LinkSessionID string
	LinkURL       string
}

// LinkStatusResult is the response from GET /api/plaid/link/status/:id.
type LinkStatusResult struct {
	Status      string
	Institution *models.Institution
}

// LinkTokenCreate starts a new Plaid Link session and returns a URL for the user to visit.
func (c *Client) LinkTokenCreate(ctx context.Context) (LinkTokenCreateResult, error) {
	// TODO: implement
	return LinkTokenCreateResult{}, nil
}

// LinkStatus polls the status of a previously created Plaid Link session.
func (c *Client) LinkStatus(ctx context.Context, sessionID string) (LinkStatusResult, error) {
	// TODO: implement
	return LinkStatusResult{}, nil
}
