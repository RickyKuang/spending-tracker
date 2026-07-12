package main

import "flag"

// runAccounts lists all accounts across the user's linked institutions.
func runAccounts(args []string) error {
	fs := flag.NewFlagSet("accounts", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement
	return nil
}
