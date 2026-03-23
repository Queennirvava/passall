package cmd

import (
	"errors"
	"fmt"
	"os"

	"passall/auth"
)

// RunAuth handles the `auth` subcommand with sub-operations: init, change.
func RunAuth(a *auth.Auth, args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: passall auth <init|change>")
		return fmt.Errorf("subcommand required")
	}
	switch args[0] {
	case "init":
		return runAuthInit(a)
	case "change":
		return runAuthChange(a)
	default:
		fmt.Fprintln(os.Stderr, "usage: passall auth <init|change>")
		return fmt.Errorf("unknown auth subcommand: %q", args[0])
	}
}

func runAuthInit(a *auth.Auth) error {
	if a.Initialized() {
		return fmt.Errorf("master password already set, use 'auth change' to update it")
	}
	pw, err := readPassword("Set master password: ")
	if err != nil {
		return err
	}
	if pw == "" {
		return fmt.Errorf("master password must not be empty")
	}
	confirm, err := readPassword("Confirm master password: ")
	if err != nil {
		return err
	}
	if pw != confirm {
		return fmt.Errorf("passwords do not match")
	}
	if err := a.Init(pw); err != nil {
		return err
	}
	fmt.Println("Master password set.")
	return nil
}

func runAuthChange(a *auth.Auth) error {
	if !a.Initialized() {
		return fmt.Errorf("master password not set, use 'auth init' first")
	}
	old, err := readPassword("Current master password: ")
	if err != nil {
		return err
	}
	newPw, err := readPassword("New master password: ")
	if err != nil {
		return err
	}
	if newPw == "" {
		return fmt.Errorf("master password must not be empty")
	}
	confirm, err := readPassword("Confirm new master password: ")
	if err != nil {
		return err
	}
	if newPw != confirm {
		return fmt.Errorf("passwords do not match")
	}
	if err := a.Change(old, newPw); err != nil {
		if errors.Is(err, auth.ErrWrongPassword) {
			return fmt.Errorf("incorrect current password")
		}
		return err
	}
	fmt.Println("Master password updated.")
	return nil
}

// RequireMasterPassword verifies the master password for protected commands.
// If not yet initialized, it prompts the user to set one first.
func RequireMasterPassword(a *auth.Auth) error {
	if !a.Initialized() {
		fmt.Println("Master password not set. Please set one to use this command.")
		if err := runAuthInit(a); err != nil {
			return err
		}
		return nil
	}
	pw, err := readPassword("Master password: ")
	if err != nil {
		return err
	}
	if err := a.Verify(pw); err != nil {
		if errors.Is(err, auth.ErrWrongPassword) {
			return fmt.Errorf("incorrect master password")
		}
		return err
	}
	return nil
}
