package cmd

import (
	"flag"
	"fmt"
	"os"

	"passall/vault"
)

func RunAdd(v *vault.Vault, args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	service := fs.String("service", "", "service name (required)")
	account := fs.String("account", "", "account name (required)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *service == "" || *account == "" {
		fmt.Fprintln(os.Stderr, "usage: passall add --service <service> --account <account>")
		return fmt.Errorf("--service and --account are required")
	}
	pw, err := readPassword("Password: ")
	if err != nil {
		return err
	}
	if pw == "" {
		return fmt.Errorf("password must not be empty")
	}
	err = v.Save(vault.Entry{
		Service:  *service,
		Account:  *account,
		Password: pw,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Added %s / %s\n", *service, *account)
	return nil
}
