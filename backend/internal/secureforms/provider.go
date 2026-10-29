package secureforms

import (
	"context"
	"errors"
)

// KeyProvider is a boundary for a future external key management service.
// Implementations must authorize organization access independently and
// avoid logging or persisting plaintext keys.
type KeyProvider interface {
	ResolveKey(ctx context.Context, organizationID, keyID string) ([]byte, error)
}

// EncryptWithProvider requires a trusted, server-selected key ID. The
// organization and submission context must come from authorized server state.
func EncryptWithProvider(ctx context.Context, provider KeyProvider, keyID string, scope Context, plaintext []byte) (Envelope, error) {
	if provider == nil || keyID == "" || scope.OrganizationID == "" {
		return Envelope{}, errors.New("missing key provider or context")
	}
	key, err := provider.ResolveKey(ctx, scope.OrganizationID, keyID)
	if err != nil { return Envelope{}, err }
	defer clear(key)
	return Encrypt(key, keyID, scope, plaintext)
}

// DecryptWithProvider resolves the envelope key in the authorized tenant
// context. The caller must not trust tenant IDs from unverified client input.
func DecryptWithProvider(ctx context.Context, provider KeyProvider, scope Context, envelope Envelope) ([]byte, error) {
	if provider == nil || scope.OrganizationID == "" || envelope.KeyID == "" {
		return nil, errors.New("missing key provider or context")
	}
	key, err := provider.ResolveKey(ctx, scope.OrganizationID, envelope.KeyID)
	if err != nil { return nil, err }
	defer clear(key)
	return Decrypt(key, scope, envelope)
}

// MemoryKeyProvider adapts the development-only Keyring to the interface.
// It must never be selected in production deployments.
type MemoryKeyProvider struct { Ring *Keyring }

func (p MemoryKeyProvider) ResolveKey(_ context.Context, organizationID, keyID string) ([]byte, error) {
	if p.Ring == nil { return nil, errors.New("keyring unavailable") }
	return p.Ring.Resolve(organizationID, keyID)
}
