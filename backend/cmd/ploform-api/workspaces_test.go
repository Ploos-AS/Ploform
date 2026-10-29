package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ploos-AS/Ploform/backend/internal/authz"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
	"github.com/Ploos-AS/Ploform/backend/internal/store"
)

type workspaceTestSessions struct{}

func (workspaceTestSessions) Lookup(context.Context, string) (session.Identity, error) {
	return session.Identity{UserID: "user-1"}, nil
}

type workspaceTestMemberships struct {
	role authz.Role
}

func (m workspaceTestMemberships) Resolve(context.Context, string, string) (authz.Role, error) {
	return m.role, nil
}

func TestWorkspaceRouteRejectsUnauthenticated(t *testing.T) {
	a := authz.Authorizer{Sessions: workspaceTestSessions{}, Memberships: workspaceTestMemberships{role: authz.Owner}}
	mux := http.NewServeMux()
	mux.Handle("GET /api/organizations/{organizationID}/workspaces", workspaceRoutes(a, store.Store{}))
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/00000000-0000-0000-0000-000000000001/workspaces", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", rec.Code)
	}
}

func TestWorkspaceRouteRejectsViewerWithoutPermission(t *testing.T) {
	a := authz.Authorizer{Sessions: workspaceTestSessions{}, Memberships: workspaceTestMemberships{role: authz.Role("unknown")}}
	mux := http.NewServeMux()
	mux.Handle("GET /api/organizations/{organizationID}/workspaces", workspaceRoutes(a, store.Store{}))
	req := httptest.NewRequest(http.MethodGet, "/api/organizations/00000000-0000-0000-0000-000000000001/workspaces", nil)
	req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "test"})
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status %d, want 403", rec.Code)
	}
}
