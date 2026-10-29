package main

import (
 "database/sql"
 "net/http"
 "net/http/httptest"
 "testing"

 _ "github.com/jackc/pgx/v5/stdlib"
)

// Routes backed by session and membership checks must reject anonymous
// requests before attempting any database query.
func TestAdminReadRoutesRejectAnonymousRequests(t *testing.T) {
 db,err:=sql.Open("pgx","postgres://localhost/unused")
 if err!=nil {t.Fatal(err)}
 defer db.Close()
 h:=handlerWithDatabases(db,db)
 org:="00000000-0000-0000-0000-000000000001"
 ws:="00000000-0000-0000-0000-000000000002"
 form:="00000000-0000-0000-0000-000000000003"
 routes:=[]string{
  "/app/"+org+"/forms",
  "/inbox/"+org+"/"+form,
  "/api/organizations/"+org+"/workspaces",
  "/api/organizations/"+org+"/workspaces/"+ws+"/forms",
  "/api/organizations/"+org+"/forms/"+form+"/submissions",
  "/api/organizations/"+org+"/forms/"+form+"/versions",
 }
 for _,path:=range routes {
  t.Run(path,func(t *testing.T){
   rec:=httptest.NewRecorder()
   h.ServeHTTP(rec,httptest.NewRequest(http.MethodGet,path,nil))
   if rec.Code!=http.StatusUnauthorized {t.Fatalf("GET %s: got %d, want 401",path,rec.Code)}
   if body:=rec.Body.String();len(body)>0 && (containsSensitiveHTML(body)) {t.Fatalf("GET %s: leaked admin HTML",path)}
  })
 }
}

func containsSensitiveHTML(body string) bool {
 return len(body)>20 && (contains(body,"<html") || contains(body,"<!doctype"))
}
func contains(haystack,needle string) bool {
 for i:=0;i+len(needle)<=len(haystack);i++ {
  if haystack[i:i+len(needle)]==needle{return true}
 }
 return false
}
