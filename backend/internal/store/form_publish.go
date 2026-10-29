package store

import (
 "context"
 "database/sql"
 "errors"
)

var ErrFormNotFound = errors.New("form not found or has no version")

// PublishForm changes only a tenant-owned draft with a validated version.
// The database transaction holds a row lock during the transition.
func (s Store) PublishForm(ctx context.Context, orgID, formID string) (Form, error) {
 var f Form
 if !validUUID(formID) { return f, ErrInvalidForm }
 err := s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  var status string
  if err := tx.QueryRowContext(ctx, "SELECT status FROM forms WHERE organization_id=$1 AND id=$2 FOR UPDATE", orgID, formID).Scan(&status); err != nil { return err }
  if status != "draft" { return ErrFormNotDraft }
  var count int
  if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM form_versions WHERE organization_id=$1 AND form_id=$2", orgID, formID).Scan(&count); err != nil { return err }
  if count == 0 { return ErrFormNotFound }
  return tx.QueryRowContext(ctx, "UPDATE forms SET status='published' WHERE organization_id=$1 AND id=$2 RETURNING id, organization_id, workspace_id, slug, title, status", orgID, formID).
   Scan(&f.ID, &f.OrganizationID, &f.WorkspaceID, &f.Slug, &f.Title, &f.Status)
 })
 if errors.Is(err, sql.ErrNoRows) { return Form{}, ErrFormNotFound }
 if err != nil { return Form{}, err }
 return f, nil
}
