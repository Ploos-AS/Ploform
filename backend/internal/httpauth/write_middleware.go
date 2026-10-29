package httpauth

import (
 "context"
 "errors"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/authz"
 "github.com/Ploos-AS/Ploform/backend/internal/session"
)

// WriteMiddleware authorizes POST requests using server-side membership.
// Callers must additionally enforce origin and CSRF checks.
func WriteMiddleware(a authz.Authorizer, permission authz.Permission, next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
   http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
   return
  }
  cookie, err := r.Cookie(session.CookieName)
  if err != nil {
   http.Error(w, "authentication required", http.StatusUnauthorized)
   return
  }
  orgID := r.PathValue("organizationID")
  if orgID == "" {
   http.Error(w, "organization required", http.StatusBadRequest)
   return
  }
  principal, err := a.Authorize(r.Context(), cookie.Value, orgID, permission)
  if err != nil {
   status := http.StatusForbidden
   if errors.Is(err, authz.ErrUnauthenticated) { status = http.StatusUnauthorized }
   http.Error(w, http.StatusText(status), status)
   return
  }
  next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, principal)))
 })
}
