package main

import (
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func TestPublishRequiresPrincipal(t *testing.T) {
 w := httptest.NewRecorder()
 r := httptest.NewRequest(http.MethodPost, "/api/organizations/x/forms/y/publish", nil)
 publishFormHandler(store.Store{}).ServeHTTP(w, r)
 if w.Code != http.StatusUnauthorized { t.Fatalf("status %d, want 401", w.Code) }
}

func TestPublishRouteFailsClosed(t *testing.T) {
 w := httptest.NewRecorder()
 r := httptest.NewRequest(http.MethodPost, "/api/organizations/00000000-0000-0000-0000-000000000001/forms/00000000-0000-0000-0000-000000000002/publish", nil)
 handler().ServeHTTP(w, r)
 if w.Code != http.StatusNotFound { t.Fatalf("status %d, want 404", w.Code) }
}
