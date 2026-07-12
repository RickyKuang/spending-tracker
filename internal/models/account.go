package models

import "time"

// Account is a bank account belonging to an Institution.
type Account struct {
	ID               string    `json:"id" db:"id"`
	InstitutionID    string    `json:"institution_id" db:"institution_id"`
	PlaidAccountID   string    `json:"plaid_account_id" db:"plaid_account_id"`
	Name             string    `json:"name" db:"name"`
	OfficialName     *string   `json:"official_name,omitempty" db:"official_name"`
	Mask             *string   `json:"mask,omitempty" db:"mask"`
	Type             string    `json:"type" db:"type"`
	Subtype          *string   `json:"subtype,omitempty" db:"subtype"`
	CurrentBalance   *float64  `json:"current_balance,omitempty" db:"current_balance"`
	AvailableBalance *float64  `json:"available_balance,omitempty" db:"available_balance"`
	ISOCurrencyCode  string    `json:"iso_currency_code" db:"iso_currency_code"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
