package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotInitialized is returned when the master password has not been set up yet.
var ErrNotInitialized = errors.New("master password not initialized")

// ErrWrongPassword is returned when the provided password does not match the stored hash.
var ErrWrongPassword = errors.New("wrong password")

// Auth manages the master password hash stored on disk.
type Auth struct {
	hashFile string
}

// New creates an Auth using the given hash file path (must be an absolute path).
func New(hashFile string) *Auth {
	return &Auth{hashFile: hashFile}
}

// Initialized reports whether the master password has been set up.
func (a *Auth) Initialized() bool {
	_, err := os.Stat(a.hashFile)
	return err == nil
}

// Init sets the master password for the first time.
// Callers should check Initialized() before calling Init.
func (a *Auth) Init(password string) error {
	return a.writeHash(password)
}

// Verify checks the provided password against the stored hash.
// Returns ErrNotInitialized if no hash file exists, ErrWrongPassword if mismatch.
func (a *Auth) Verify(password string) error {
	data, err := os.ReadFile(a.hashFile)
	if errors.Is(err, os.ErrNotExist) {
		return ErrNotInitialized
	}
	if err != nil {
		return fmt.Errorf("auth: read hash file: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword(data, []byte(password)); err != nil {
		return ErrWrongPassword
	}
	return nil
}

// Change updates the master password after verifying the old one.
func (a *Auth) Change(oldPassword, newPassword string) error {
	if err := a.Verify(oldPassword); err != nil {
		return err
	}
	return a.writeHash(newPassword)
}

// writeHash hashes password with bcrypt and writes it to the hash file.
func (a *Auth) writeHash(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("auth: hash password: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(a.hashFile), 0700); err != nil {
		return fmt.Errorf("auth: create dir: %w", err)
	}
	if err := os.WriteFile(a.hashFile, hash, 0600); err != nil {
		return fmt.Errorf("auth: write hash file: %w", err)
	}
	return nil
}
