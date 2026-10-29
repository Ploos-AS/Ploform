package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Ploos-AS/Ploform/backend/internal/secureforms"
)

// EncryptedSubmissionInput is an internal service input, not an HTTP payload.
// OrganizationID must be resolved from an authorized principal.
type EncryptedSubmissionInput struct {
	OrganizationID string
	FormID string
	FormVersionID string
	SubmissionID string
	KeyID string
	Payload []byte
}

// SaveSensitiveSubmission encrypts before beginning the tenant transaction.
// It never writes plaintext to response_json or encrypted_payload.
func (s Store) SaveSensitiveSubmission(ctx context.Context, keys secureforms.KeyProvider, input EncryptedSubmissionInput) error {
	if !validUUID(input.OrganizationID) || !validUUID(input.FormID) ||
		!validUUID(input.FormVersionID) || !validUUID(input.SubmissionID) {
		return errors.New("invalid submission identifiers")
	}
	if keys == nil || input.KeyID == "" {
		return errors.New("encryption key provider required")
	}
	scope := secureforms.Context{
		OrganizationID: input.OrganizationID,
		FormID: input.FormID,
		FormVersionID: input.FormVersionID,
		SubmissionID: input.SubmissionID,
	}
	envelope, err := secureforms.EncryptWithProvider(ctx, keys, input.KeyID, scope, input.Payload)
	if err != nil { return fmt.Errorf("encrypt submission: %w", err) }
	encoded, err := json.Marshal(envelope)
	if err != nil { return fmt.Errorf("serialize envelope: %w", err) }
	return s.InTenant(ctx, input.OrganizationID, func(tx *sql.Tx) error {
		// The composite join proves that the form version belongs to the
		// form in the authenticated encryption context.
		result, err := tx.ExecContext(ctx, `
			INSERT INTO submissions
				(id, organization_id, form_version_id, response_json, encryption_mode, encrypted_payload)
			SELECT $1, $2, v.id, '{}'::jsonb, 'sensitive', $4::jsonb
			FROM form_versions v
			WHERE v.id = $3 AND v.organization_id = $2 AND v.form_id = $5
		`, input.SubmissionID, input.OrganizationID, input.FormVersionID, string(encoded), input.FormID)
		if err != nil { return err }
		n, err := result.RowsAffected()
		if err != nil { return err }
		if n != 1 { return errors.New("form version not found in organization") }
		return nil
	})
}
