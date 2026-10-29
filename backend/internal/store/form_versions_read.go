package store

import (
 "context"
 "database/sql"
 "encoding/json"

)

func (s Store) ListFormVersions(ctx context.Context, orgID, formID string) ([]FormVersion, error) {
 result := []FormVersion{}
 if !validUUID(formID) { return result, sql.ErrNoRows }
 err := s.InTenant(ctx, orgID, func(tx *sql.Tx) error {
  rows, err := tx.QueryContext(ctx, "SELECT v.id, v.organization_id, v.form_id, v.version, v.schema_json FROM form_versions v JOIN forms f ON f.id=v.form_id AND f.organization_id=v.organization_id WHERE v.organization_id=$1 AND v.form_id=$2 ORDER BY v.version DESC", orgID, formID)
  if err != nil { return err }
  defer rows.Close()
  for rows.Next() {
   var item FormVersion
   var raw []byte
   if err := rows.Scan(&item.ID, &item.OrganizationID, &item.FormID, &item.Version, &raw); err != nil { return err }
   if err := json.Unmarshal(raw, &item.Schema); err != nil { return err }
   result = append(result, item)
  }
  return rows.Err()
 })
 return result, err
}
