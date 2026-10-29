package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ploos-AS/Ploform/backend/internal/authz"
	"github.com/Ploos-AS/Ploform/backend/internal/httpauth"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
	"github.com/Ploos-AS/Ploform/backend/internal/store"
)

func handler() http.Handler { return handlerWithDB(nil) }

func handlerWithDB(db *sql.DB) http.Handler { return handlerWithDatabases(db, nil) }

func handlerWithDatabases(db, authDB *sql.DB) http.Handler {
	mux := http.NewServeMux()
	if db != nil && authDB != nil {
		sessions := session.Repository{DB: authDB}
		a := authz.Authorizer{Sessions: sessions, Memberships: authz.MembershipRepository{DB: authDB}}
		mux.HandleFunc("GET /api/me", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			cookie, err := r.Cookie(session.CookieName)
			if err != nil { http.Error(w, "authentication required", http.StatusUnauthorized); return }
			identity, err := sessions.Lookup(r.Context(), cookie.Value)
			if err != nil { http.Error(w, "authentication required", http.StatusUnauthorized); return }
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"user_id": identity.UserID})
		})
		mux.Handle("GET /api/organizations/{organizationID}/workspaces", workspaceRoutes(a, store.Store{DB: db}))
		mux.Handle("GET /api/organizations/{organizationID}/workspaces/{workspaceID}/forms", formsListRoute(a, store.Store{DB: db}))
		mux.Handle("GET /api/csrf", csrfTokenHandler(sessions))
		mux.Handle("POST /api/logout", httpauth.RequireSameOrigin(requireCSRF(logoutHandler(sessions))))
	}
	mux.Handle("GET /readyz", readinessHandler(db))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	return mux
}

func main() {
	if err := httpauth.ValidateTrustedProxyCIDRs(os.Getenv("PLOFORM_TRUSTED_PROXY_IPS")); err != nil {
		log.Fatalf("trusted proxy configuration: %v", err)
	}
	addr := os.Getenv("PLOFORM_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	// No administrative endpoints exist yet. Fail closed when a database
	// is configured but unavailable, rather than starting a degraded service.
	var db *sql.DB
	if os.Getenv("PLOFORM_DATABASE_URL") != "" || os.Getenv("PLOFORM_APP_PASSWORD_FILE") != "" {
		var err error
		db, err = connectDatabase(context.Background())
		if err != nil {
			log.Fatalf("database startup check failed: %v", err)
		}
		defer db.Close()
	}
	var authDB *sql.DB
	if os.Getenv("PLOFORM_AUTH_DATABASE_URL") != "" || os.Getenv("PLOFORM_AUTH_PASSWORD_FILE") != "" {
		if db == nil {
			log.Fatal("authentication requires PLOFORM_DATABASE_URL")
		}
		var err error
		authDB, err = connectAuthDatabase(context.Background())
		if err != nil {
			log.Fatalf("authentication database startup check failed: %v", err)
		}
		defer authDB.Close()
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           handlerWithDatabases(db, authDB),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Ploform development API listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
