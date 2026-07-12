package main

import (
	"flag"
	"fmt"
)

// runAuth dispatches to the auth signup/signin/signout subcommands.
func runAuth(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: spender auth <signup|signin|signout> [args]")
	}

	switch args[0] {
	case "signup":
		return runAuthSignUp(args[1:])
	case "signin":
		return runAuthSignIn(args[1:])
	case "signout":
		return runAuthSignOut(args[1:])
	default:
		return fmt.Errorf("usage: spender auth <signup|signin|signout> [args]")
	}
}

func runAuthSignUp(args []string) error {
	fs := flag.NewFlagSet("auth signup", flag.ExitOnError)
	email := fs.String("email", "", "account email")
	password := fs.String("password", "", "account password")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _ = *email, *password

	// TODO: implement
	return nil
}

func runAuthSignIn(args []string) error {
	fs := flag.NewFlagSet("auth signin", flag.ExitOnError)
	email := fs.String("email", "", "account email")
	password := fs.String("password", "", "account password")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, _ = *email, *password

	// TODO: implement
	return nil
}

func runAuthSignOut(args []string) error {
	fs := flag.NewFlagSet("auth signout", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	// TODO: implement
	return nil
}
