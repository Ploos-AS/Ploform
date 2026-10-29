package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

func TestContactTemplateEndpoint(t *testing.T) {
 req := httptest.NewRequest(http.MethodGet, "/api/templates/contact", nil)
 rec := httptest.NewRecorder()
 handler().ServeHTTP(rec, req)
 if rec.Code != http.StatusOK { t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String()) }
 if got := rec.Header().Get("Cache-Control"); got != "no-store" { t.Fatalf("cache control = %q", got) }
 var response struct {
  ID string `json:"id"`
  Name string `json:"name"`
  Schema formschema.Schema `json:"schema"`
 }
 if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil { t.Fatal(err) }
 if response.ID != "contact" || response.Name == "" { t.Fatalf("unexpected template: %+v", response) }
 if err := formschema.Validate(response.Schema); err != nil { t.Fatalf("invalid template schema: %v", err) }
 if len(response.Schema.Fields) != 4 { t.Fatalf("expected 4 fields, got %d", len(response.Schema.Fields)) }
}

func TestContactTemplateIsReadOnly(t *testing.T) {
 req := httptest.NewRequest(http.MethodPost, "/api/templates/contact", nil)
 rec := httptest.NewRecorder()
 handler().ServeHTTP(rec, req)
 if rec.Code != http.StatusMethodNotAllowed { t.Fatalf("status = %d, want 405", rec.Code) }
}
