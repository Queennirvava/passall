package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"passall/vault"
)

func RunDelete(v *vault.Vault, args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	service := fs.String("service", "", "service name (required)")
	account := fs.String("account", "", "account name (required)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *service == "" || *account == "" {
		fmt.Fprintln(os.Stderr, "usage: passall delete --service <service> --account <account>")
		return fmt.Errorf("--service and --account are required")
	}
	err := v.Delete(*service, *account)
	if err != nil {
		if errors.Is(err, vault.ErrEntryNotFound) {
			return fmt.Errorf("no entry found for %s / %s", *service, *account)
		}
		return err
	}
	fmt.Printf("Deleted %s / %s\n", *service, *account)
	return nil
}
