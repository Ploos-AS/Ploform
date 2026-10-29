package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Ploos-AS/Ploform/backend/internal/httpauth"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type csrfSessionLookup interface {
	Lookup(context.Context, string) (session.Identity, error)
}

// csrfTokenHandler issues a token only for an active, database-verified session.
func csrfTokenHandler(s csrfSessionLookup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Vary", "Cookie")
		cookie, err := r.Cookie(session.CookieName)
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if _, err := s.Lookup(r.Context(), cookie.Value); err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		token, err := httpauth.NewCSRFToken(cookie.Value)
		if err != nil {
			http.Error(w, "token unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"csrf_token": token})
	})
}

// requireCSRF validates a session-bound token before calling a write handler.
// The caller must additionally enforce same-origin HTTPS and authorization.
func requireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(session.CookieName)
		if err != nil || !httpauth.ValidCSRFToken(cookie.Value, r.Header.Get("X-CSRF-Token")) {
			http.Error(w, "invalid CSRF token", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
