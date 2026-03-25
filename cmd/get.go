package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"passall/auth"
	"passall/vault"
)

func RunGet(v *vault.Vault, a *auth.Auth, args []string) error {
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

	// Prompt for master password. Empty input = skip, show nothing.
	pw, err := readPassword("Master password (Enter to skip): ")
	if err != nil {
		return err
	}
	showPassword := false
	if pw != "" {
		if err := a.Verify(pw); err != nil {
			if errors.Is(err, auth.ErrWrongPassword) {
				return fmt.Errorf("incorrect master password")
			}
			return err
		}
		showPassword = true
	}

	if len(entries) == 1 && *account != "" {
		e := entries[0]
		if showPassword {
			fmt.Printf("service:  %s\naccount:  %s\npassword: %s\n", e.Service, e.Account, e.Password)
		} else {
			fmt.Printf("service:  %s\naccount:  %s\n", e.Service, e.Account)
		}
		return nil
	}
	if showPassword {
		fmt.Printf("%-20s %-20s %s\n", "SERVICE", "ACCOUNT", "PASSWORD")
		for _, e := range entries {
			fmt.Printf("%-20s %-20s %s\n", e.Service, e.Account, e.Password)
		}
	} else {
		fmt.Printf("%-20s %s\n", "SERVICE", "ACCOUNT")
		for _, e := range entries {
			fmt.Printf("%-20s %s\n", e.Service, e.Account)
		}
	}
	return nil
}
