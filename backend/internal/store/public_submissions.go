package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
 "github.com/Ploos-AS/Ploform/backend/internal/submission"
)

// SavePublicSubmission validates against the current published version and
// stores the normalized response within the same tenant-scoped transaction.
// Callers must enforce request-size limits and abuse controls before invoking it.
func (s Store) SavePublicSubmission(ctx context.Context, orgID, workspaceSlug, formSlug string, values map[string][]string) (string,error) {
 if !validUUID(orgID)||ValidateSlug(workspaceSlug)!=nil||ValidateSlug(formSlug)!=nil {return "",ErrFormNotFound}
 var id string
 err:=s.InTenant(ctx,orgID,func(tx *sql.Tx)error{
  var versionID string
  var raw []byte
  err:=tx.QueryRowContext(ctx,`SELECT v.id,v.schema_json
FROM forms f JOIN workspaces w ON w.id=f.workspace_id AND w.organization_id=f.organization_id
JOIN LATERAL (SELECT id,schema_json FROM form_versions
 WHERE organization_id=f.organization_id AND form_id=f.id ORDER BY version DESC LIMIT 1) v ON true
WHERE f.organization_id=$1 AND w.slug=$2 AND f.slug=$3 AND f.status='published'
FOR SHARE OF f`,orgID,workspaceSlug,formSlug).Scan(&versionID,&raw)
  if errors.Is(err,sql.ErrNoRows){return ErrFormNotFound}
  if err!=nil{return err}
  var schema formschema.Schema
  if err=json.Unmarshal(raw,&schema);err!=nil{return err}
  normalized,err:=submission.Validate(schema,values)
  if err!=nil{return err}
  encoded,err:=json.Marshal(normalized)
  if err!=nil{return err}
  return tx.QueryRowContext(ctx,`INSERT INTO submissions (organization_id,form_version_id,response_json)
VALUES ($1,$2,$3::jsonb) RETURNING id`,orgID,versionID,string(encoded)).Scan(&id)
 })
 if err!=nil{return "",err}
 return id,nil
}
