package main

import (
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func TestPublicFormUnavailableWithoutDatabase(t *testing.T) {
 r := httptest.NewRequest(http.MethodGet, "/api/public/organizations/00000000-0000-0000-0000-000000000001/workspaces/demo/forms/contact", nil)
 w := httptest.NewRecorder()
 handler().ServeHTTP(w,r)
 if w.Code != http.StatusNotFound { t.Fatalf("status=%d want 404",w.Code) }
}

func TestPublicFormInvalidIdentifiers(t *testing.T) {
 r := httptest.NewRequest(http.MethodGet, "/api/public/organizations/not-a-uuid/workspaces/demo/forms/contact", nil)
 r.SetPathValue("organizationID","not-a-uuid")
 r.SetPathValue("workspaceSlug","demo")
 r.SetPathValue("formSlug","contact")
 w := httptest.NewRecorder()
 publicFormHandler(store.Store{}).ServeHTTP(w,r)
 if w.Code != http.StatusNotFound { t.Fatalf("status=%d want 404",w.Code) }
}
