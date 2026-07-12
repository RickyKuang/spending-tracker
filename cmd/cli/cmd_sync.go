package main

import "flag"

// runSync triggers an on-demand sync and prints a per-institution summary.
func runSync(args []string) error {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement
	return nil
}
