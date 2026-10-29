package main

import (
 "encoding/json"
 "errors"
 "io"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func createFormDraftHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  principal, ok := httpauth.Principal(r.Context())
  if !ok { http.Error(w, "authentication required", http.StatusUnauthorized); return }
  if r.Header.Get("Content-Type") != "application/json" {
   http.Error(w, "application/json required", http.StatusUnsupportedMediaType); return
  }
  var input struct {
   Slug string `json:"slug"`
   Title string `json:"title"`
  }
  dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
  dec.DisallowUnknownFields()
  if err := dec.Decode(&input); err != nil {
   http.Error(w, "invalid JSON", http.StatusBadRequest); return
  }
  var extra any
  if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
   http.Error(w, "unexpected JSON content", http.StatusBadRequest); return
  }
  form, err := s.CreateFormDraft(r.Context(), principal.OrganizationID, r.PathValue("workspaceID"), input.Slug, input.Title)
  if err != nil {
   if errors.Is(err, store.ErrInvalidForm) {
    http.Error(w, "invalid form draft or workspace", http.StatusBadRequest); return
   }
   http.Error(w, "database unavailable", http.StatusServiceUnavailable); return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  w.WriteHeader(http.StatusCreated)
  _ = json.NewEncoder(w).Encode(form)
 })
}
