package main

import (
 "encoding/json"
 "errors"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func publishFormHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  principal, ok := httpauth.Principal(r.Context())
  if !ok { http.Error(w, "authentication required", http.StatusUnauthorized); return }
  form, err := s.PublishForm(r.Context(), principal.OrganizationID, r.PathValue("formID"))
  if err != nil {
   if errors.Is(err, store.ErrInvalidForm) { http.Error(w, "invalid form id", http.StatusBadRequest); return }
   if errors.Is(err, store.ErrFormNotFound) { http.Error(w, "form not found or missing version", http.StatusNotFound); return }
   if errors.Is(err, store.ErrFormNotDraft) { http.Error(w, "form is not a draft", http.StatusConflict); return }
   http.Error(w, "database unavailable", http.StatusServiceUnavailable); return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  w.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(w).Encode(form)
 })
}
