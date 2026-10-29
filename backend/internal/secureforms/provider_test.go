package secureforms

import (
	"bytes"
	"context"
	"testing"
)

func TestProviderRoundtripAndIsolation(t *testing.T) {
	var ring Keyring
	if err := ring.Register("org-a", "key-a", bytes.Repeat([]byte{0x71}, KeySize)); err != nil { t.Fatal(err) }
	provider := MemoryKeyProvider{Ring: &ring}
	scope := Context{OrganizationID: "org-a", FormID: "form", FormVersionID: "v1", SubmissionID: "sub"}
	envelope, err := EncryptWithProvider(context.Background(), provider, "key-a", scope, []byte("answer"))
	if err != nil { t.Fatal(err) }
	plain, err := DecryptWithProvider(context.Background(), provider, scope, envelope)
	if err != nil || string(plain) != "answer" { t.Fatalf("roundtrip: %v", err) }
	other := scope
	other.OrganizationID = "org-b"
	if _, err := DecryptWithProvider(context.Background(), provider, other, envelope); err == nil {
		t.Fatal("cross-tenant key resolution succeeded")
	}
}

func TestProviderRejectsNil(t *testing.T) {
	scope := Context{OrganizationID: "org", FormID: "form", FormVersionID: "v", SubmissionID: "sub"}
	if _, err := EncryptWithProvider(context.Background(), nil, "key", scope, nil); err == nil {
		t.Fatal("nil provider accepted")
	}
	if _, err := DecryptWithProvider(context.Background(), nil, scope, Envelope{KeyID: "key"}); err == nil {
		t.Fatal("nil provider accepted")
	}
}
