package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type sessionRevoker interface {
	Revoke(context.Context, string) error
}

// logoutHandler invalidates the server-side session before clearing its cookie.
// Only mount this behind a verified HTTPS same-origin boundary.
func logoutHandler(s sessionRevoker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		cookie, err := r.Cookie(session.CookieName)
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		if s == nil {
			http.Error(w, "session service unavailable", http.StatusServiceUnavailable)
			return
		}
		err = s.Revoke(r.Context(), cookie.Value)
		if err != nil && !errors.Is(err, session.ErrNotFound) {
			http.Error(w, "session service unavailable", http.StatusServiceUnavailable)
			return
		}
		session.ClearCookie(w)
		w.WriteHeader(http.StatusNoContent)
	})
}
