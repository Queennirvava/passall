package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"passall/vault"
)

func RunGet(v *vault.Vault, args []string) error {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	service := fs.String("service", "", "service name (required)")
	account := fs.String("account", "", "account name (optional)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *service == "" {
		fmt.Fprintln(os.Stderr, "usage: passall get --service <service> [--account <account>]")
		return fmt.Errorf("--service is required")
	}
	entries, err := v.Get(*service, *account)
	if err != nil {
		if errors.Is(err, vault.ErrEntryNotFound) {
			return fmt.Errorf("no entry found for %s / %s", *service, *account)
		}
		return err
	}
	if len(entries) == 1 && *account != "" {
		e := entries[0]
		fmt.Printf("service:  %s\naccount:  %s\npassword: %s\n", e.Service, e.Account, e.Password)
		return nil
	}
	fmt.Printf("%-20s %s\n", "SERVICE", "ACCOUNT")
	for _, e := range entries {
		fmt.Printf("%-20s %s\n", e.Service, e.Account)
	}
	return nil
}
