package main

import (
    "database/sql"
    "net/http"
    "net/http/httptest"
    "testing"

    _ "github.com/jackc/pgx/v5/stdlib"
)

// A partial or missing database configuration must never expose form or
// authentication routes. These checks do not require a running PostgreSQL.
func TestSensitiveRoutesFailClosedWithoutBothDatabases(t *testing.T) {
    db, err := sql.Open("pgx", "postgres://localhost/unused")
    if err != nil { t.Fatal(err) }
    defer db.Close()

    org := "00000000-0000-0000-0000-000000000001"
    workspace := "00000000-0000-0000-0000-000000000002"
    form := "00000000-0000-0000-0000-000000000003"
    routes := []struct {
        method, path string
    }{
        {"GET", "/api/me"},
        {"GET", "/api/csrf"},
        {"POST", "/api/logout"},
        {"GET", "/api/organizations/" + org + "/workspaces"},
        {"GET", "/api/organizations/" + org + "/workspaces/" + workspace + "/forms"},
        {"POST", "/api/organizations/" + org + "/workspaces/" + workspace + "/forms"},
        {"GET", "/api/organizations/" + org + "/forms/" + form + "/versions"},
        {"GET", "/api/organizations/" + org + "/forms/" + form + "/submissions"},
        {"GET", "/app/" + org + "/forms"},
        {"GET", "/inbox/" + org + "/" + form},
        {"POST", "/api/organizations/" + org + "/forms/" + form + "/versions"},
    }
    configurations := []struct {
        name string
        app, auth *sql.DB
    }{
        {"neither", nil, nil},
        {"app-only", db, nil},
        {"auth-only", nil, db},
    }
    for _, cfg := range configurations {
        t.Run(cfg.name, func(t *testing.T) {
            h := handlerWithDatabases(cfg.app, cfg.auth)
            for _, route := range routes {
                t.Run(route.method+" "+route.path, func(t *testing.T) {
                    rec := httptest.NewRecorder()
                    req := httptest.NewRequest(route.method, route.path, nil)
                    h.ServeHTTP(rec, req)
                    if rec.Code != http.StatusNotFound {
                        t.Fatalf("%s %s: got %d, want 404", route.method, route.path, rec.Code)
                    }
                })
            }
        })
    }
}
