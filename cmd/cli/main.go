// Command spender is the Go reference CLI for the spending-tracker backend API.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "auth":
		err = runAuth(os.Args[2:])
	case "link":
		err = runLink(os.Args[2:])
	case "sync":
		err = runSync(os.Args[2:])
	case "accounts":
		err = runAccounts(os.Args[2:])
	case "transactions":
		err = runTransactions(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "spender:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: spender <auth|link|sync|accounts|transactions> [args]")
}
