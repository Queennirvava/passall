package main

import (
	"fmt"
	"os"

	"passall/auth"
	"passall/cmd"
	"passall/config"
	"passall/vault"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	v, err := vault.NewVault(cfg.Storage.VaultDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	a := auth.New(cfg.Auth.MasterPasswordHashFile)

	subcmd := os.Args[1]
	args := os.Args[2:]

	var runErr error
	switch subcmd {
	case "add":
		runErr = cmd.RunAdd(v, args)
	case "get":
		runErr = cmd.RunGet(v, args)
	case "update":
		if runErr = cmd.RequireMasterPassword(a); runErr == nil {
			runErr = cmd.RunUpdate(v, args)
		}
	case "delete":
		if runErr = cmd.RequireMasterPassword(a); runErr == nil {
			runErr = cmd.RunDelete(v, args)
		}
	case "list":
		runErr = cmd.RunList(v, args)
	case "auth":
		runErr = cmd.RunAuth(a, args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %q\n", subcmd)
		printUsage()
		os.Exit(1)
	}

	if runErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", runErr)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `usage: passall <command> [options]

Commands:
  add     --service <s> --account <a>            add a new entry
  get     --service <s> [--account <a>]          get entry/entries
  update  --service <s> --account <a>            update password
  delete  --service <s> --account <a>            delete an entry
  list    [--service <s>]                        list entries
  auth    <init|change>                          manage master password`)
}
