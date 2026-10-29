package main

import (
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "strings"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func formVersionsListHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  p, ok := httpauth.Principal(r.Context())
  if !ok { http.Error(w, "authentication required", 401); return }
  versions, err := s.ListFormVersions(r.Context(), p.OrganizationID, r.PathValue("formID"))
  if err != nil { http.Error(w, "versions unavailable", 503); return }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  _ = json.NewEncoder(w).Encode(map[string]any{"versions":versions})
 })
}

func formVersionsCreateHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  p, ok := httpauth.Principal(r.Context())
  if !ok { http.Error(w, "authentication required", 401); return }
  if !strings.EqualFold(strings.TrimSpace(strings.Split(r.Header.Get("Content-Type"), ";")[0]), "application/json") {
   http.Error(w, "application/json required", 415); return
  }
  dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 65536))
  dec.DisallowUnknownFields()
  var schema formschema.Schema
  if err := dec.Decode(&schema); err != nil { http.Error(w, "invalid JSON", 400); return }
  var extra any
  if err := dec.Decode(&extra); !errors.Is(err, io.EOF) { http.Error(w, "unexpected JSON", 400); return }
  version, err := s.AddFormVersion(r.Context(), p.OrganizationID, r.PathValue("formID"), schema)
  if err != nil {
   if errors.Is(err, formschema.ErrInvalid) { http.Error(w, "invalid schema", 400); return }
   if errors.Is(err, store.ErrFormNotDraft) { http.Error(w, "form is not a draft", 409); return }
   http.Error(w, "version creation unavailable", 503); return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  w.WriteHeader(http.StatusCreated)
  _ = json.NewEncoder(w).Encode(version)
 })
}
