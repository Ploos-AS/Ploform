package main

import (
 "encoding/json"
 "errors"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func publicFormHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  form, err := s.GetPublicForm(r.Context(), r.PathValue("organizationID"), r.PathValue("workspaceSlug"), r.PathValue("formSlug"))
  if err != nil {
   if errors.Is(err, store.ErrFormNotFound) { http.Error(w, "form not found", http.StatusNotFound); return }
   http.Error(w, "form unavailable", http.StatusServiceUnavailable); return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  _ = json.NewEncoder(w).Encode(form)
 })
}
