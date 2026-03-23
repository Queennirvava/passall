package cmd

import (
	"flag"
	"fmt"
	"os"

	"passall/vault"
)

func RunList(v *vault.Vault, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	service := fs.String("service", "", "filter by service (optional)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	var entries []vault.Entry
	var err error
	if *service != "" {
		entries, err = v.List(*service)
	} else {
		entries, err = v.ListAll()
	}
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "no entries found")
		return nil
	}
	fmt.Printf("%-20s %s\n", "SERVICE", "ACCOUNT")
	for _, e := range entries {
		fmt.Printf("%-20s %s\n", e.Service, e.Account)
	}
	return nil
}
