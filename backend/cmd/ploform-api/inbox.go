package main

import (
 "encoding/json"
 "net/http"
 "strconv"

 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func inboxListRoute(a authz.Authorizer,s store.Store) http.Handler {
 return httpauth.Middleware(a,authz.ReadResponses,func(r *http.Request)string{return r.PathValue("organizationID")},http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  principal,ok:=httpauth.Principal(r.Context())
  if !ok {http.Error(w,"authentication required",http.StatusUnauthorized);return}
  limit:=25
  offset:=0
  if raw:=r.URL.Query().Get("limit");raw!=""{
   n,err:=strconv.Atoi(raw)
   if err!=nil||n<1||n>100{http.Error(w,"invalid limit",http.StatusBadRequest);return}
   limit=n
  }
  if raw:=r.URL.Query().Get("offset");raw!=""{
   n,err:=strconv.Atoi(raw)
   if err!=nil||n<0||n>1000000{http.Error(w,"invalid offset",http.StatusBadRequest);return}
   offset=n
  }
  items,err:=s.ListFormSubmissions(r.Context(),principal.OrganizationID,r.PathValue("formID"),limit,offset)
  if err!=nil{http.Error(w,"inbox unavailable",http.StatusServiceUnavailable);return}
  w.Header().Set("Content-Type","application/json; charset=utf-8")
  _=json.NewEncoder(w).Encode(map[string]any{"submissions":items,"limit":limit,"offset":offset})
 }))
}
