package models

import "time"

// Transaction is a single bank transaction synced from Plaid.
// Plaid amount convention: positive = money OUT (debit), negative = money IN.
type Transaction struct {
	ID                   string     `json:"id" db:"id"`
	AccountID            string     `json:"account_id" db:"account_id"`
	PlaidTransactionID   string     `json:"plaid_transaction_id" db:"plaid_transaction_id"`
	Amount               float64    `json:"amount" db:"amount"`
	ISOCurrencyCode      string     `json:"iso_currency_code" db:"iso_currency_code"`
	PostedDate           time.Time  `json:"posted_date" db:"posted_date"`
	AuthorizedDate       *time.Time `json:"authorized_date,omitempty" db:"authorized_date"`
	Name                 string     `json:"name" db:"name"`
	MerchantName         *string    `json:"merchant_name,omitempty" db:"merchant_name"`
	CategoryPrimary      *string    `json:"category_primary,omitempty" db:"category_primary"`
	CategoryDetailed     *string    `json:"category_detailed,omitempty" db:"category_detailed"`
	Pending              bool       `json:"pending" db:"pending"`
	PendingTransactionID *string    `json:"pending_transaction_id,omitempty" db:"pending_transaction_id"`
	PaymentChannel       *string    `json:"payment_channel,omitempty" db:"payment_channel"`
	RemovedAt            *time.Time `json:"removed_at,omitempty" db:"removed_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}
