// Package session supplies cryptographically random opaque session tokens.
// Persistence and authentication are intentionally separate.
package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"
)

const tokenBytes = 32

var ErrExpired = errors.New("session expired")

// NewToken returns a bearer secret suitable for an HttpOnly, Secure cookie.
// Store only its SHA-256 digest, never the raw token.
func NewToken() (string, [32]byte, error) {
	var raw [tokenBytes]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", [32]byte{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(raw[:])
	return token, sha256.Sum256(raw[:]), nil
}

// Digest decodes a token and computes its storage key. Invalid tokens fail
// closed before any database lookup.
func Digest(token string) ([32]byte, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != tokenBytes || base64.RawURLEncoding.EncodeToString(raw) != token {
		return [32]byte{}, errors.New("invalid session token")
	}
	return sha256.Sum256(raw), nil
}

// Active rejects expired or revoked sessions; expiry is exclusive.
func Active(now, expiresAt time.Time, revoked bool) error {
	if revoked || !now.Before(expiresAt) {
		return ErrExpired
	}
	return nil
}
