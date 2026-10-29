package secureforms

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestEncryptDecryptAndTenantIsolation(t *testing.T) {
	key := bytes.Repeat([]byte{0x42}, KeySize)
	ctx := Context{"org-a", "form-a", "version-a", "submission-a"}
	message := []byte("confidential response")
	env, err := Encrypt(key, "key-1", ctx, message)
	if err != nil { t.Fatal(err) }
	plain, err := Decrypt(key, ctx, env)
	if err != nil || !bytes.Equal(plain, message) { t.Fatalf("roundtrip failed: %v", err) }
	other := ctx
	other.OrganizationID = "org-b"
	if _, err := Decrypt(key, other, env); err == nil { t.Fatal("cross-tenant decryption succeeded") }
	other = ctx
	other.SubmissionID = "submission-b"
	if _, err := Decrypt(key, other, env); err == nil { t.Fatal("cross-submission decryption succeeded") }
	env.Ciphertext = base64.RawURLEncoding.EncodeToString([]byte("tampered"))
	if _, err := Decrypt(key, ctx, env); err == nil { t.Fatal("tampered ciphertext accepted") }
}

func TestEncryptionUsesFreshNonce(t *testing.T) {
	key := bytes.Repeat([]byte{0x13}, KeySize)
	ctx := Context{"org", "form", "version", "submission"}
	a, err := Encrypt(key, "key-1", ctx, []byte("same"))
	if err != nil { t.Fatal(err) }
	b, err := Encrypt(key, "key-1", ctx, []byte("same"))
	if err != nil { t.Fatal(err) }
	if a.Nonce == b.Nonce || a.Ciphertext == b.Ciphertext { t.Fatal("repeated encryption output") }
}

func TestRejectsMissingContextAndKey(t *testing.T) {
	key := bytes.Repeat([]byte{0x13}, KeySize)
	if _, err := Encrypt(key, "key", Context{}, []byte("x")); err == nil { t.Fatal("missing context accepted") }
	if _, err := Encrypt([]byte("short"), "key", Context{"o","f","v","s"}, []byte("x")); err == nil { t.Fatal("short key accepted") }
}
