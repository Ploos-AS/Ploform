package secureforms

import (
	"crypto/subtle"
	"errors"
	"sync"
)

// Keyring is a development-only in-memory key resolver. It deliberately
// provides no persistence, recovery, rotation or production key custody.
// Production deployments must replace it with an external key manager.
type Keyring struct {
	mu sync.RWMutex
	keys map[string]keyEntry
}

type keyEntry struct {
	organizationID string
	key []byte
}

// Register accepts externally provisioned key material and copies it.
// A key identifier is immutable and belongs to exactly one organization.
func (r *Keyring) Register(organizationID, keyID string, key []byte) error {
	if organizationID == "" || keyID == "" || len(key) != KeySize {
		return errors.New("invalid key registration")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.keys == nil { r.keys = make(map[string]keyEntry) }
	if _, exists := r.keys[keyID]; exists {
		return errors.New("key identifier already registered")
	}
	r.keys[keyID] = keyEntry{organizationID: organizationID, key: append([]byte(nil), key...)}
	return nil
}

// Resolve only releases key material to callers naming the correct
// organization. Authorization of the caller must happen before this call.
func (r *Keyring) Resolve(organizationID, keyID string) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, exists := r.keys[keyID]
	if !exists || organizationID == "" ||
		subtle.ConstantTimeCompare([]byte(entry.organizationID), []byte(organizationID)) != 1 {
		return nil, errors.New("key unavailable")
	}
	return append([]byte(nil), entry.key...), nil
}
