package main

import (
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

// Malformed or unauthenticated requests must not reach the database.
func TestCreateContactFormRequiresPrincipal(t *testing.T) {
 h := createContactFormHandler(store.Store{})
 r := httptest.NewRequest(http.MethodPost, "/api/organizations/x/workspaces/y/contact-forms", strings.NewReader(`{"slug":"contact","title":"Contact"}`))
 r.Header.Set("Content-Type", "application/json")
 w := httptest.NewRecorder()
 h.ServeHTTP(w, r)
 if w.Code != http.StatusUnauthorized { t.Fatalf("got %d, want 401", w.Code) }
}

func TestCreateContactFormRouteRequiresAuthentication(t *testing.T) {
 // Both handles are non-nil, but intentionally never connected: an
 // anonymous request must be rejected before any SQL is attempted.
 // Authentication is checked by the existing middleware.
 // The no-database case must not register the sensitive route.
 path := "/api/organizations/00000000-0000-0000-0000-000000000001/workspaces/00000000-0000-0000-0000-000000000002/contact-forms"
 r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"slug":"contact","title":"Contact"}`))
 r.Header.Set("Content-Type", "application/json")
 w := httptest.NewRecorder()
 handler().ServeHTTP(w, r)
 if w.Code != http.StatusNotFound { t.Fatalf("got %d, want 404", w.Code) }
}
