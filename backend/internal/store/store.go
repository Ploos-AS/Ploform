// Package store provides tenant-scoped PostgreSQL operations.
package store

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "regexp"
 "strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{1,62}$`)

type Store struct { DB *sql.DB }

type Workspace struct { ID, OrganizationID, Slug, Name string }

func ValidateSlug(slug string) error {
 if !slugPattern.MatchString(slug) { return errors.New("slug must be 2–63 lowercase letters, digits or hyphens") }
 return nil
}

// InTenant runs all tenant operations inside a transaction with a transaction-local
// tenant setting. orgID MUST come from verified authorization, never user input alone.
func (s Store) InTenant(ctx context.Context, orgID string, fn func(*sql.Tx) error) (err error) {
 if s.DB == nil { return errors.New("database not configured") }
 if !validUUID(orgID) { return errors.New("invalid organization id") }
 tx, err := s.DB.BeginTx(ctx, nil)
 if err != nil { return fmt.Errorf("begin tenant transaction: %w", err) }
 defer tx.Rollback()
 if _, err = tx.ExecContext(ctx, "SELECT set_config('app.organization_id', $1, true)", orgID); err != nil {
  return fmt.Errorf("set tenant context: %w", err)
 }
 if err = fn(tx); err != nil { return err }
 return tx.Commit()
}

func (s Store) CreateWorkspace(ctx context.Context, orgID, slug, name string) (Workspace, error) {
 var w Workspace
 if err := ValidateSlug(slug); err != nil { return w, err }
 name = strings.TrimSpace(name)
 if name == "" { return w, errors.New("workspace name is required") }
 err := s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  return tx.QueryRowContext(ctx,
   "INSERT INTO workspaces (organization_id, slug, name) VALUES ($1, $2, $3) RETURNING id, organization_id, slug, name",
   orgID, slug, name).Scan(&w.ID, &w.OrganizationID, &w.Slug, &w.Name)
 })
 return w, err
}

func (s Store) ListWorkspaces(ctx context.Context, orgID string) ([]Workspace, error) {
 result := []Workspace{}
 err := s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  rows, err := tx.QueryContext(ctx,
   "SELECT id, organization_id, slug, name FROM workspaces WHERE organization_id = $1 ORDER BY slug", orgID)
  if err != nil { return err }
  defer rows.Close()
  for rows.Next() {
   var w Workspace
   if err := rows.Scan(&w.ID, &w.OrganizationID, &w.Slug, &w.Name); err != nil { return err }
   result = append(result, w)
  }
  return rows.Err()
 })
 return result, err
}

func validUUID(v string) bool {
 if len(v) != 36 { return false }
 for i, c := range v {
  if i == 8 || i == 13 || i == 18 || i == 23 {
   if c != '-' { return false }
  } else if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) { return false }
 }
 return true
}
