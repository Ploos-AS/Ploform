package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "strings"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

// CreateContactForm creates a draft and its first version atomically.
func (s Store) CreateContactForm(ctx context.Context, orgID, workspaceID, slug, title string) (Form, FormVersion, error) {
 var form Form
 var version FormVersion
 if !validUUID(workspaceID) || ValidateSlug(slug) != nil || strings.TrimSpace(title) == "" || len(title) > 200 {
  return form, version, ErrInvalidForm
 }
 schema := formschema.Contact()
 encoded, err := json.Marshal(schema)
 if err != nil { return form, version, err }
 err = s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  err := tx.QueryRowContext(ctx, `INSERT INTO forms (organization_id, workspace_id, slug, title)
   VALUES ($1,$2,$3,$4) RETURNING id, organization_id, workspace_id, slug, title, status`,
   orgID, workspaceID, slug, strings.TrimSpace(title)).
   Scan(&form.ID, &form.OrganizationID, &form.WorkspaceID, &form.Slug, &form.Title, &form.Status)
  if err != nil { return err }
  return tx.QueryRowContext(ctx, `INSERT INTO form_versions (organization_id, form_id, version, schema_json)
   VALUES ($1,$2,1,$3::jsonb) RETURNING id, organization_id, form_id, version`,
   orgID, form.ID, string(encoded)).
   Scan(&version.ID, &version.OrganizationID, &version.FormID, &version.Version)
 })
 if err != nil { return Form{}, FormVersion{}, err }
 version.Schema = schema
 return form, version, nil
}
