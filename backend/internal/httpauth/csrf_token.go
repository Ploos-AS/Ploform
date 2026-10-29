package httpauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

const csrfNonceBytes = 32

// NewCSRFToken binds a random nonce to an opaque session token using HMAC.
// The secret is derived from the session and never stored in client-readable
// cookies. A fresh nonce makes tokens non-deterministic.
func NewCSRFToken(sessionToken string) (string, error) {
	if sessionToken == "" {
		return "", errors.New("missing session")
	}
	nonce := make([]byte, csrfNonceBytes)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	mac := csrfMAC(sessionToken, nonce)
	raw := append(nonce, mac...)
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func ValidCSRFToken(sessionToken, token string) bool {
	if sessionToken == "" {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != csrfNonceBytes+sha256.Size || base64.RawURLEncoding.EncodeToString(raw) != token {
		return false
	}
	want := csrfMAC(sessionToken, raw[:csrfNonceBytes])
	return subtle.ConstantTimeCompare(raw[csrfNonceBytes:], want) == 1
}

func csrfMAC(sessionToken string, nonce []byte) []byte {
	mac := hmac.New(sha256.New, []byte(sessionToken))
	mac.Write([]byte("ploform-csrf-v1"))
	mac.Write(nonce)
	return mac.Sum(nil)
}
