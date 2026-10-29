package store

import (
 "context"
 "database/sql"
)

// Form is the tenant-scoped metadata returned by the draft forms API.
type Form struct {
 ID string `json:"id"`
 OrganizationID string `json:"organization_id"`
 WorkspaceID string `json:"workspace_id"`
 Slug string `json:"slug"`
 Title string `json:"title"`
 Status string `json:"status"`
}

// ListForms returns only forms in the authorized organization and workspace.
// RLS is enforced by InTenant as an independent database-side boundary.
func (s Store) ListForms(ctx context.Context, organizationID, workspaceID string) ([]Form, error) {
 result := []Form{}
 if !validUUID(workspaceID) { return result, sql.ErrNoRows }
 err := s.InTenant(ctx, organizationID, func(tx *sql.Tx) error {
  rows, err := tx.QueryContext(ctx,
   "SELECT id, organization_id, workspace_id, slug, title, status FROM forms WHERE organization_id = $1 AND workspace_id = $2 ORDER BY slug",
   organizationID, workspaceID)
  if err != nil { return err }
  defer rows.Close()
  for rows.Next() {
   var f Form
   if err := rows.Scan(&f.ID, &f.OrganizationID, &f.WorkspaceID, &f.Slug, &f.Title, &f.Status); err != nil { return err }
   result = append(result, f)
  }
  return rows.Err()
 })
 return result, err
}
