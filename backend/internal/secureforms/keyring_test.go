package secureforms

import (
	"bytes"
	"testing"
)

func TestKeyringTenantIsolation(t *testing.T) {
	var ring Keyring
	key := bytes.Repeat([]byte{0x55}, KeySize)
	if err := ring.Register("org-a", "key-1", key); err != nil { t.Fatal(err) }
	key[0] = 0
	resolved, err := ring.Resolve("org-a", "key-1")
	if err != nil || resolved[0] != 0x55 { t.Fatalf("key registration was not copied: %v", err) }
	resolved[0] = 0
	again, err := ring.Resolve("org-a", "key-1")
	if err != nil || again[0] != 0x55 { t.Fatalf("resolved key was not copied: %v", err) }
	if _, err := ring.Resolve("org-b", "key-1"); err == nil { t.Fatal("cross-tenant key access") }
	if _, err := ring.Resolve("org-a", "missing"); err == nil { t.Fatal("unknown key resolved") }
	if err := ring.Register("org-b", "key-1", bytes.Repeat([]byte{1}, KeySize)); err == nil { t.Fatal("key ID reused") }
}

func TestKeyringRejectsInvalidRegistration(t *testing.T) {
	var ring Keyring
	for _, tc := range []struct{ org, id string; key []byte }{
		{"", "key", bytes.Repeat([]byte{1}, KeySize)},
		{"org", "", bytes.Repeat([]byte{1}, KeySize)},
		{"org", "key", []byte("short")},
	} {
		if err := ring.Register(tc.org, tc.id, tc.key); err == nil { t.Fatal("invalid registration accepted") }
	}
}
