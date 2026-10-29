// Package secureforms provides an authenticated-encryption building block.
// Key custody and HTTP submission integration are deliberately out of scope.
package secureforms

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	Version = 1
	KeySize = 32
	nonceSize = 12
)

type Context struct {
	OrganizationID string `json:"organization_id"`
	FormID string `json:"form_id"`
	FormVersionID string `json:"form_version_id"`
	SubmissionID string `json:"submission_id"`
}

type Envelope struct {
	Version int `json:"version"`
	Algorithm string `json:"algorithm"`
	KeyID string `json:"key_id"`
	Nonce string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func (c Context) associatedData() ([]byte, error) {
	if c.OrganizationID == "" || c.FormID == "" || c.FormVersionID == "" || c.SubmissionID == "" {
		return nil, errors.New("missing encryption context")
	}
	return json.Marshal(c)
}

func Encrypt(key []byte, keyID string, context Context, plaintext []byte) (Envelope, error) {
	if len(key) != KeySize || keyID == "" {
		return Envelope{}, errors.New("invalid encryption key or key identifier")
	}
	aad, err := context.associatedData()
	if err != nil { return Envelope{}, err }
	block, err := aes.NewCipher(key)
	if err != nil { return Envelope{}, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return Envelope{}, err }
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil { return Envelope{}, fmt.Errorf("nonce generation: %w", err) }
	sealed := gcm.Seal(nil, nonce, plaintext, aad)
	return Envelope{
		Version: Version, Algorithm: "AES-256-GCM", KeyID: keyID,
		Nonce: base64.RawURLEncoding.EncodeToString(nonce),
		Ciphertext: base64.RawURLEncoding.EncodeToString(sealed),
	}, nil
}

func Decrypt(key []byte, context Context, envelope Envelope) ([]byte, error) {
	if len(key) != KeySize || envelope.Version != Version || envelope.Algorithm != "AES-256-GCM" || envelope.KeyID == "" {
		return nil, errors.New("unsupported envelope or invalid key")
	}
	aad, err := context.associatedData()
	if err != nil { return nil, err }
	nonce, err := base64.RawURLEncoding.DecodeString(envelope.Nonce)
	if err != nil || len(nonce) != nonceSize { return nil, errors.New("invalid nonce") }
	ciphertext, err := base64.RawURLEncoding.DecodeString(envelope.Ciphertext)
	if err != nil { return nil, errors.New("invalid ciphertext") }
	block, err := aes.NewCipher(key)
	if err != nil { return nil, err }
	gcm, err := cipher.NewGCM(block)
	if err != nil { return nil, err }
	plaintext, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil { return nil, errors.New("authentication failed") }
	return plaintext, nil
}
