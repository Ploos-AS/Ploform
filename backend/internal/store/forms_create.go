package store

import (
 "context"
 "database/sql"
 "errors"
 "strings"
)

var ErrInvalidForm = errors.New("invalid form draft")

func (s Store) CreateFormDraft(ctx context.Context, orgID, workspaceID, slug, title string) (Form, error) {
 var f Form
 if !validUUID(workspaceID) || ValidateSlug(slug) != nil { return f, ErrInvalidForm }
 title = strings.TrimSpace(title)
 if title == "" || len(title) > 200 { return f, ErrInvalidForm }
 err := s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  return tx.QueryRowContext(ctx,
   "INSERT INTO forms (organization_id, workspace_id, slug, title, status) SELECT $1, id, $3, $4, 'draft' FROM workspaces WHERE organization_id = $1 AND id = $2 RETURNING id, organization_id, workspace_id, slug, title, status",
   orgID, workspaceID, slug, title).Scan(&f.ID, &f.OrganizationID, &f.WorkspaceID, &f.Slug, &f.Title, &f.Status)
 })
 if errors.Is(err, sql.ErrNoRows) { return Form{}, ErrInvalidForm }
 return f, err
}
