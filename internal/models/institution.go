package models

import "time"

// Institution is one linked bank / Plaid Item.
type Institution struct {
	ID                  string     `json:"id" db:"id"`
	UserID              string     `json:"user_id" db:"user_id"`
	PlaidItemID         string     `json:"plaid_item_id" db:"plaid_item_id"`
	PlaidAccessTokenEnc []byte     `json:"-" db:"plaid_access_token_enc"`
	PlaidInstitutionID  string     `json:"plaid_institution_id" db:"plaid_institution_id"`
	Name                string     `json:"name" db:"name"`
	SyncCursor          *string    `json:"sync_cursor,omitempty" db:"sync_cursor"`
	LastSyncedAt        *time.Time `json:"last_synced_at,omitempty" db:"last_synced_at"`
	SyncStatus          string     `json:"sync_status" db:"sync_status"`
	SyncError           *string    `json:"sync_error,omitempty" db:"sync_error"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}
