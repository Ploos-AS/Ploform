package httpauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ploos-AS/Ploform/backend/internal/authz"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type fakeSessions struct{}
func (fakeSessions) Lookup(context.Context, string) (session.Identity, error) {
	return session.Identity{UserID: "user-1"}, nil
}
type fakeMembers struct{ role authz.Role }
func (f fakeMembers) Resolve(context.Context, string, string) (authz.Role, error) {
	return f.role, nil
}

func TestMiddleware(t *testing.T) {
	for _, tc := range []struct {
		name   string
		method string
		cookie bool
		role   authz.Role
		want   int
	}{
		{"authorized", http.MethodGet, true, authz.Viewer, http.StatusOK},
		{"no cookie", http.MethodGet, false, authz.Viewer, http.StatusUnauthorized},
		{"insufficient role", http.MethodGet, true, authz.Role("invalid"), http.StatusForbidden},
		{"writes disabled", http.MethodPost, true, authz.Owner, http.StatusMethodNotAllowed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := authz.Authorizer{Sessions: fakeSessions{}, Memberships: fakeMembers{role: tc.role}}
			h := Middleware(a, authz.ReadResponses, func(*http.Request) string { return "org-1" }, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				p, ok := Principal(r.Context())
				if !ok || p.UserID != "user-1" || p.OrganizationID != "org-1" {
					t.Fatal("verified principal missing")
				}
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest(tc.method, "/test", nil)
			if tc.cookie {
				req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fake"})
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("got status %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
