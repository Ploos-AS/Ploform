package main

import (
 "encoding/json"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/httpauth"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

// formsListRoute is read-only until authenticated CSRF-protected writes exist.
func formsListRoute(a authz.Authorizer, s store.Store) http.Handler {
 return httpauth.Middleware(a, authz.ReadResponses, func(r *http.Request) string {
  return r.PathValue("organizationID")
 }, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  principal, ok := httpauth.Principal(r.Context())
  if !ok { http.Error(w, "authentication required", http.StatusUnauthorized); return }
  items, err := s.ListForms(r.Context(), principal.OrganizationID, r.PathValue("workspaceID"))
  if err != nil { http.Error(w, "form listing unavailable", http.StatusServiceUnavailable); return }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Cache-Control", "no-store")
  _ = json.NewEncoder(w).Encode(map[string]any{"forms": items})
 }))
}
