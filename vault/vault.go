package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrEntryExists is returned when attempting to save an entry that already exists.
var ErrEntryExists = errors.New("entry already exists")

// ErrEntryNotFound is returned when an entry cannot be found.
var ErrEntryNotFound = errors.New("entry not found")

// Vault manages password entries stored as JSON files on disk.
// Each service has its own file: <Path>/<service>.json
type Vault struct {
	Path string
}

// NewVault creates a Vault rooted at path.
// The directory is created if it does not already exist.
func NewVault(path string) (*Vault, error) {
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("failed to create vault dir: %w", err)
	}
	return &Vault{Path: path}, nil
}

// readEntries reads all entries for a service. Returns an empty slice if the
// file does not exist.
func (v *Vault) readEntries(service string) ([]Entry, error) {
	data, err := os.ReadFile(filepath.Join(v.Path, service+".json"))
	if errors.Is(err, os.ErrNotExist) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading entries %q: %w", service, err)
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("reading entries %q: %w", service, err)
	}
	return entries, nil
}

// writeEntries marshals entries and writes them to the service's JSON file.
func (v *Vault) writeEntries(service string, entries []Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("writing entries %q: %w", service, err)
	}
	if err := os.WriteFile(filepath.Join(v.Path, service+".json"), data, 0600); err != nil {
		return fmt.Errorf("writing entries %q: %w", service, err)
	}
	return nil
}

// Save stores a new entry. Returns ErrEntryExists if service+account already exists.
func (v *Vault) Save(entry Entry) error {
	entries, err := v.readEntries(entry.Service)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Account == entry.Account {
			return ErrEntryExists
		}
	}
	now := time.Now().Format(time.RFC3339)
	entry.CreatedAt = now
	entry.UpdatedAt = now
	entries = append(entries, entry)
	return v.writeEntries(entry.Service, entries)
}

// Get retrieves entries by service and optional account.
// If account is empty, all entries for the service are returned.
// If account is specified, returns the single matching entry.
// Returns ErrEntryNotFound if no matching entry exists.
// service must not be empty.
func (v *Vault) Get(service, account string) ([]Entry, error) {
	if service == "" {
		return nil, fmt.Errorf("vault: service must not be empty")
	}
	entries, err := v.readEntries(service)
	if err != nil {
		return nil, err
	}
	if account == "" {
		return entries, nil
	}
	for i := range entries {
		if entries[i].Account == account {
			return []Entry{entries[i]}, nil
		}
	}
	return nil, ErrEntryNotFound
}

// Update replaces an existing entry (matched by service+account).
// Only the password and UpdatedAt fields are updated.
// Returns ErrEntryNotFound if no matching entry exists.
func (v *Vault) Update(entry Entry) error {
	entries, err := v.readEntries(entry.Service)
	if err != nil {
		return err
	}
	for i := range entries {
		if entries[i].Account == entry.Account {
			entries[i].Password = entry.Password
			entries[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return v.writeEntries(entry.Service, entries)
		}
	}
	return ErrEntryNotFound
}

// Delete removes an entry identified by service and account.
// Returns ErrEntryNotFound if no matching entry exists.
func (v *Vault) Delete(service, account string) error {
	entries, err := v.readEntries(service)
	if err != nil {
		return err
	}
	for i, e := range entries {
		if e.Account == account {
			entries = append(entries[:i], entries[i+1:]...)
			return v.writeEntries(service, entries)
		}
	}
	return ErrEntryNotFound
}

// List returns all entries for the given service.
func (v *Vault) List(service string) ([]Entry, error) {
	return v.readEntries(service)
}

// ListAll returns every entry across all services in the vault.
func (v *Vault) ListAll() ([]Entry, error) {
	pattern := filepath.Join(v.Path, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("vault: list all: %w", err)
	}
	var all []Entry
	for _, f := range files {
		service := f[len(v.Path)+1 : len(f)-len(".json")]
		entries, err := v.readEntries(service)
		if err != nil {
			return nil, err
		}
		all = append(all, entries...)
	}
	return all, nil
}
