package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ploos-AS/Ploform/backend/internal/store"
)

func TestWorkspaceJSONContract(t *testing.T) {
	w := store.Workspace{
		ID: "workspace-id",
		OrganizationID: "organization-id",
		Slug: "course-registration",
		Name: "Course registration",
	}
	data, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var obj map[string]string
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatal(err)
	}
	for key, expected := range map[string]string{
		"id": "workspace-id",
		"organization_id": "organization-id",
		"slug": "course-registration",
		"name": "Course registration",
	} {
		if got := obj[key]; got != expected {
			t.Errorf("%s = %q, want %q", key, got, expected)
		}
	}
	if len(obj) != 4 {
		t.Errorf("unexpected fields: %v", obj)
	}
}

func TestWorkspaceRouteDisabledWithoutBothDatabases(t *testing.T) {
	path := "/api/organizations/00000000-0000-0000-0000-000000000001/workspaces"
	for _, h := range []http.Handler{
		handlerWithDatabases(nil, nil),
		handlerWithDatabases(nil, nil),
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status %d, want 404 when databases are not configured", rec.Code)
		}
	}
}
