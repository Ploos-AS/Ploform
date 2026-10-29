package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

var ErrFormNotDraft = errors.New("form is not an editable draft")

type FormVersion struct {
 ID string `json:"id"`
 OrganizationID string `json:"organization_id"`
 FormID string `json:"form_id"`
 Version int `json:"version"`
 Schema formschema.Schema `json:"schema"`
}

// AddFormVersion atomically allocates the next version for a draft form.
// The form row lock serializes concurrent writers.
func (s Store) AddFormVersion(ctx context.Context, orgID, formID string, schema formschema.Schema) (FormVersion, error) {
 var result FormVersion
 if !validUUID(formID) || formschema.Validate(schema) != nil { return result, formschema.ErrInvalid }
 encoded, err := json.Marshal(schema)
 if err != nil { return result, err }
 err = s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  var status string
  err := tx.QueryRowContext(ctx, "SELECT status FROM forms WHERE organization_id=$1 AND id=$2 FOR UPDATE", orgID, formID).Scan(&status)
  if err != nil { return err }
  if status != "draft" { return ErrFormNotDraft }
  return tx.QueryRowContext(ctx, `INSERT INTO form_versions (organization_id, form_id, version, schema_json)
SELECT $1, $2, COALESCE(MAX(version),0)+1, $3::jsonb
FROM form_versions WHERE organization_id=$1 AND form_id=$2
RETURNING id, organization_id, form_id, version`, orgID, formID, string(encoded)).
   Scan(&result.ID, &result.OrganizationID, &result.FormID, &result.Version)
 })
 if err != nil { return FormVersion{}, err }
 result.Schema = schema
 return result, nil
}
