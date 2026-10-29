package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Ploos-AS/Ploform/backend/internal/secureforms"
)

// ReadSensitiveSubmission is internal-only. The organization ID must come
// from an authorized principal with permission to read form responses.
func (s Store) ReadSensitiveSubmission(ctx context.Context, keys secureforms.KeyProvider, organizationID, submissionID string) ([]byte, error) {
	if !validUUID(organizationID) || !validUUID(submissionID) {
		return nil, errors.New("invalid submission identifiers")
	}
	if keys == nil {
		return nil, errors.New("encryption key provider required")
	}
	var envelopeJSON []byte
	var scope secureforms.Context
	err := s.InTenant(ctx, organizationID, func(tx *sql.Tx) error {
		return tx.QueryRowContext(ctx, `
			SELECT s.organization_id, v.form_id, s.form_version_id, s.id, s.encrypted_payload
			FROM submissions s
			JOIN form_versions v
			  ON v.id = s.form_version_id AND v.organization_id = s.organization_id
			WHERE s.id = $1 AND s.organization_id = $2
			  AND s.encryption_mode = 'sensitive'
			  AND s.response_json = '{}'::jsonb
		`, submissionID, organizationID).Scan(
			&scope.OrganizationID, &scope.FormID, &scope.FormVersionID,
			&scope.SubmissionID, &envelopeJSON,
		)
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("encrypted submission not found")
	}
	if err != nil {
		return nil, fmt.Errorf("read encrypted submission: %w", err)
	}
	var envelope secureforms.Envelope
	if err := json.Unmarshal(envelopeJSON, &envelope); err != nil {
		return nil, errors.New("invalid encrypted submission envelope")
	}
	plaintext, err := secureforms.DecryptWithProvider(ctx, keys, scope, envelope)
	if err != nil {
		return nil, fmt.Errorf("decrypt submission: %w", err)
	}
	return plaintext, nil
}
