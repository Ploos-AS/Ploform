package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

// PublicForm contains only fields intended for anonymous visitors.
type PublicForm struct {
 Slug string `json:"slug"`
 Title string `json:"title"`
 Schema formschema.Schema `json:"schema"`
}

// GetPublicForm reads a published form using the explicitly supplied organization
// as a tenant context. This is NOT authorization: the SQL published predicate and
// deliberately narrow response are the public access boundary.
func (s Store) GetPublicForm(ctx context.Context, orgID, workspaceSlug, formSlug string) (PublicForm, error) {
 var result PublicForm
 if !validUUID(orgID) || ValidateSlug(workspaceSlug) != nil || ValidateSlug(formSlug) != nil {
  return result, ErrFormNotFound
 }
 if s.DB == nil { return result, errors.New("database not configured") }
 tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{ReadOnly:true})
 if err != nil { return result, err }
 defer tx.Rollback()
 if _, err = tx.ExecContext(ctx, "SELECT set_config('app.organization_id', $1, true)", orgID); err != nil { return result, err }
 var raw []byte
 err = tx.QueryRowContext(ctx, `SELECT f.slug, f.title, v.schema_json
 FROM forms f JOIN workspaces w ON w.id=f.workspace_id AND w.organization_id=f.organization_id
 JOIN LATERAL (
  SELECT schema_json FROM form_versions
  WHERE organization_id=f.organization_id AND form_id=f.id
  ORDER BY version DESC LIMIT 1
 ) v ON true
 WHERE f.organization_id=$1 AND w.slug=$2 AND f.slug=$3 AND f.status='published'`,
 orgID, workspaceSlug, formSlug).Scan(&result.Slug, &result.Title, &raw)
 if errors.Is(err, sql.ErrNoRows) { return PublicForm{}, ErrFormNotFound }
 if err != nil { return PublicForm{}, err }
 if err = json.Unmarshal(raw, &result.Schema); err != nil { return PublicForm{}, err }
 if err = formschema.Validate(result.Schema); err != nil { return PublicForm{}, err }
 if err = tx.Commit(); err != nil { return PublicForm{}, err }
 return result, nil
}
