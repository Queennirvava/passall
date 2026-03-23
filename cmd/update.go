package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"passall/vault"
)

func RunUpdate(v *vault.Vault, args []string) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	service := fs.String("service", "", "service name (required)")
	account := fs.String("account", "", "account name (required)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *service == "" || *account == "" {
		fmt.Fprintln(os.Stderr, "usage: passall update --service <service> --account <account>")
		return fmt.Errorf("--service and --account are required")
	}
	pw, err := readPassword("New password: ")
	if err != nil {
		return err
	}
	if pw == "" {
		return fmt.Errorf("password must not be empty")
	}
	err = v.Update(vault.Entry{
		Service:  *service,
		Account:  *account,
		Password: pw,
	})
	if err != nil {
		if errors.Is(err, vault.ErrEntryNotFound) {
			return fmt.Errorf("no entry found for %s / %s", *service, *account)
		}
		return err
	}
	fmt.Printf("Updated %s / %s\n", *service, *account)
	return nil
}
