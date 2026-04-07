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

	switch os.Args[1] {
	case "sync":
		fmt.Println("sync: not implemented yet")
	case "accounts":
		fmt.Println("accounts: not implemented yet")
	case "transactions":
		fmt.Println("transactions: not implemented yet")
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: spender <sync|accounts|transactions>")
}
