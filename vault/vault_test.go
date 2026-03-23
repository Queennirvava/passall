package vault

import (
	"errors"
	"testing"
)

func newTestVault(t *testing.T) *Vault {
	t.Helper()
	v, err := NewVault(t.TempDir())
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}
	return v
}

// TestSaveAndGet verifies that a saved entry can be retrieved.
func TestSaveAndGet(t *testing.T) {
	v := newTestVault(t)

	entry := Entry{
		Service:  "github",
		Account:  "alice",
		Password: "secret123",
	}
	if err := v.Save(entry); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := v.Get("github", "alice")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got[0].Password != "secret123" {
		t.Errorf("got password %q, want %q", got[0].Password, "secret123")
	}
	if got[0].CreatedAt == "" {
		t.Error("CreatedAt should be set")
	}
	if got[0].UpdatedAt == "" {
		t.Error("UpdatedAt should be set")
	}
}

// TestSaveDuplicate verifies that saving the same service+account returns ErrEntryExists.
func TestSaveDuplicate(t *testing.T) {
	v := newTestVault(t)

	entry := Entry{Service: "github", Account: "alice", Password: "pass1"}
	if err := v.Save(entry); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	entry.Password = "pass2"
	err := v.Save(entry)
	if !errors.Is(err, ErrEntryExists) {
		t.Errorf("got %v, want ErrEntryExists", err)
	}
}

// TestUpdate verifies that updating an entry changes the stored password.
func TestUpdate(t *testing.T) {
	v := newTestVault(t)

	original := Entry{Service: "github", Account: "alice", Password: "old"}
	if err := v.Save(original); err != nil {
		t.Fatalf("Save: %v", err)
	}

	updated := Entry{Service: "github", Account: "alice", Password: "new"}
	if err := v.Update(updated); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := v.Get("github", "alice")
	if err != nil {
		t.Fatalf("Get after Update: %v", err)
	}
	if got[0].Password != "new" {
		t.Errorf("got password %q, want %q", got[0].Password, "new")
	}
}

// TestDelete verifies that deleting an entry makes subsequent Get return ErrEntryNotFound.
func TestDelete(t *testing.T) {
	v := newTestVault(t)

	entry := Entry{Service: "github", Account: "alice", Password: "pass"}
	if err := v.Save(entry); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if err := v.Delete("github", "alice"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := v.Get("github", "alice")
	if !errors.Is(err, ErrEntryNotFound) {
		t.Errorf("got %v, want ErrEntryNotFound", err)
	}
}

// TestGetByServiceOnly verifies that Get with empty account returns all entries for the service.
func TestGetByServiceOnly(t *testing.T) {
	v := newTestVault(t)

	for _, e := range []Entry{
		{Service: "github", Account: "alice", Password: "p1"},
		{Service: "github", Account: "bob", Password: "p2"},
		{Service: "google", Account: "alice", Password: "p3"},
	} {
		if err := v.Save(e); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	got, err := v.Get("github", "")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d entries, want 2", len(got))
	}
}

// TestGetEmptyServiceError verifies that Get with empty service returns an error.
func TestGetEmptyServiceError(t *testing.T) {
	v := newTestVault(t)
	_, err := v.Get("", "alice")
	if err == nil {
		t.Error("expected error for empty service, got nil")
	}
}

// TestList verifies that List returns all entries for the given service.
func TestList(t *testing.T) {
	v := newTestVault(t)

	entries := []Entry{
		{Service: "github", Account: "alice", Password: "p1"},
		{Service: "github", Account: "bob", Password: "p2"},
		{Service: "google", Account: "alice", Password: "p3"},
	}
	for _, e := range entries {
		if err := v.Save(e); err != nil {
			t.Fatalf("Save %v: %v", e, err)
		}
	}

	list, err := v.List("github")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("List(github) returned %d entries, want 2", len(list))
	}
}

// TestListAll verifies that ListAll returns entries across all services.
func TestListAll(t *testing.T) {
	v := newTestVault(t)

	entries := []Entry{
		{Service: "github", Account: "alice", Password: "p1"},
		{Service: "google", Account: "bob", Password: "p2"},
		{Service: "aws", Account: "carol", Password: "p3"},
	}
	for _, e := range entries {
		if err := v.Save(e); err != nil {
			t.Fatalf("Save %v: %v", e, err)
		}
	}

	all, err := v.ListAll()
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("ListAll returned %d entries, want 3", len(all))
	}
}
