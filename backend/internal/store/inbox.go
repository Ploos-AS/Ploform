package store

import (
 "context"
 "database/sql"
 "encoding/json"
 "time"
)

// InboxItem represents a stored response for an organization's form.
// Only the requested form's submissions are returned, newest first.
type InboxItem struct {
 ID string `json:"id"`
 FormVersionID string `json:"form_version_id"`
 Response map[string]any `json:"response"`
 CreatedAt time.Time `json:"created_at"`
}

func (s Store) ListFormSubmissions(ctx context.Context, organizationID, formID string, limit, offset int) ([]InboxItem,error) {
 items:=[]InboxItem{}
 if !validUUID(formID) {return items,sql.ErrNoRows}
 if limit<1 || limit>100 {limit=25}
 if offset<0 || offset>1000000 {offset=0}
 err:=s.InTenant(ctx,organizationID,func(tx *sql.Tx) error {
  rows,err:=tx.QueryContext(ctx,`SELECT s.id,s.form_version_id,s.response_json,s.created_at
FROM submissions s
JOIN form_versions v ON v.id=s.form_version_id AND v.organization_id=s.organization_id
JOIN forms f ON f.id=v.form_id AND f.organization_id=v.organization_id
WHERE s.organization_id=$1 AND f.id=$2
ORDER BY s.created_at DESC,s.id DESC LIMIT $3 OFFSET $4`,organizationID,formID,limit,offset)
  if err!=nil{return err}
  defer rows.Close()
  for rows.Next(){
   var item InboxItem
   var raw []byte
   if err:=rows.Scan(&item.ID,&item.FormVersionID,&raw,&item.CreatedAt);err!=nil{return err}
   if err:=json.Unmarshal(raw,&item.Response);err!=nil{return err}
   items=append(items,item)
  }
  return rows.Err()
 })
 return items,err
}
