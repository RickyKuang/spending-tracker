package main

import "flag"

// runTransactions lists transactions, optionally filtered by account, date range, and pagination cursor.
func runTransactions(args []string) error {
	fs := flag.NewFlagSet("transactions", flag.ExitOnError)
	account := fs.String("account", "", "filter by account ID")
	from := fs.String("from", "", "filter by start date (YYYY-MM-DD)")
	to := fs.String("to", "", "filter by end date (YYYY-MM-DD)")
	limit := fs.Int("limit", 0, "maximum number of transactions to return")
	cursor := fs.String("cursor", "", "pagination cursor from a previous call")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _, _, _, _ = *account, *from, *to, *limit, *cursor

	// TODO: implement
	return nil
}
