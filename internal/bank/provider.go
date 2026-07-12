// Package bank defines the provider-agnostic interface for bank data sources.
package bank

import "context"

// LinkTokenCreateResult is returned when initiating a Plaid Link session.
type LinkTokenCreateResult struct {
	LinkToken  string
	Expiration string
}

// ExchangeResult is returned after exchanging a public token for a permanent access token.
type ExchangeResult struct {
	AccessToken string
	ItemID      string
}

// SyncResult holds one page of transaction changes from a provider sync call.
type SyncResult struct {
	Added      []Transaction
	Modified   []Transaction
	Removed    []string
	NextCursor string
	HasMore    bool
}

// Transaction is the provider-agnostic shape of a single transaction.
type Transaction struct {
	ProviderTransactionID string
	ProviderAccountID     string
	Amount                float64
	ISOCurrencyCode       string
	PostedDate            string
	AuthorizedDate        string
	Name                  string
	MerchantName          string
	CategoryPrimary       string
	CategoryDetailed      string
	Pending               bool
	PendingTransactionID  string
	PaymentChannel        string
}

// Account is the provider-agnostic shape of a single bank account.
type Account struct {
	ProviderAccountID string
	Name              string
	OfficialName      string
	Mask              string
	Type              string
	Subtype           string
	CurrentBalance    float64
	AvailableBalance  float64
	ISOCurrencyCode   string
}

// Provider abstracts a bank data source (Plaid, Teller, SimpleFIN, ...).
type Provider interface {
	// CreateLinkToken starts a new Link session for a given user.
	CreateLinkToken(ctx context.Context, clientUserID string) (LinkTokenCreateResult, error)

	// ExchangePublicToken exchanges a Link public token for a permanent access token.
	ExchangePublicToken(ctx context.Context, publicToken string) (ExchangeResult, error)

	// SyncTransactions fetches the next page of transaction changes for an Item.
	SyncTransactions(ctx context.Context, accessToken string, cursor string) (SyncResult, error)

	// GetAccounts fetches current account and balance data for an Item.
	GetAccounts(ctx context.Context, accessToken string) ([]Account, error)
}
