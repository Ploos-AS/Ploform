// Package httpauth protects API handlers with session and membership checks.
// Authentication issuance and CSRF enforcement must be implemented before
// any state-changing HTTP routes are enabled.
package httpauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ploos-AS/Ploform/backend/internal/authz"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type principalKey struct{}

// Principal returns a verified principal installed by Middleware.
func Principal(ctx context.Context) (authz.Principal, bool) {
	p, ok := ctx.Value(principalKey{}).(authz.Principal)
	return p, ok
}

// Middleware checks a cookie session and server-side organization membership.
// The organization ID is a requested scope, not proof of authorization.
// Only use this for read-only endpoints until CSRF protection is integrated.
func Middleware(a authz.Authorizer, permission authz.Permission, organizationID func(*http.Request) string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cookie, err := r.Cookie(session.CookieName)
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		orgID := organizationID(r)
		if orgID == "" {
			http.Error(w, "organization required", http.StatusBadRequest)
			return
		}
		p, err := a.Authorize(r.Context(), cookie.Value, orgID, permission)
		if err != nil {
			status := http.StatusForbidden
			if errors.Is(err, authz.ErrUnauthenticated) {
				status = http.StatusUnauthorized
			}
			http.Error(w, http.StatusText(status), status)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
	})
}
