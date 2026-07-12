package main

import "flag"

// runLink opens a browser to complete Plaid Link and polls the backend for completion.
func runLink(args []string) error {
	fs := flag.NewFlagSet("link", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement
	return nil
}
