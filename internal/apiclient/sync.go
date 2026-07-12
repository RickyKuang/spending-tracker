package apiclient

import "context"

// InstitutionSyncSummary summarizes the sync result for a single institution.
type InstitutionSyncSummary struct {
	ID       string
	Added    int
	Modified int
	Removed  int
	Error    string
}

// SyncResult is the response from POST /api/sync.
type SyncResult struct {
	Institutions []InstitutionSyncSummary
}

// Sync triggers an on-demand sync of all of the user's linked institutions.
func (c *Client) Sync(ctx context.Context) (SyncResult, error) {
	// TODO: implement
	return SyncResult{}, nil
}
