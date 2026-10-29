package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ploos-AS/Ploform/backend/internal/httpauth"
	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type csrfTestLookup struct{ err error }

func (s csrfTestLookup) Lookup(context.Context, string) (session.Identity, error) {
	return session.Identity{UserID: "user", ExpiresAt: time.Now().Add(time.Hour)}, s.err
}

func TestCSRFEndpointRequiresSession(t *testing.T) {
	rec := httptest.NewRecorder()
	csrfTokenHandler(csrfTestLookup{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/csrf", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", rec.Code)
	}
}

func TestCSRFEndpointRejectsRevokedSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/csrf", nil)
	req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "invalid"})
	rec := httptest.NewRecorder()
	csrfTokenHandler(csrfTestLookup{err: errors.New("revoked")}).ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status %d, want 401", rec.Code)
	}
}

func TestCSRFGuard(t *testing.T) {
	sessionToken, _, err := session.NewToken()
	if err != nil {
		t.Fatal(err)
	}
	csrfToken, err := httpauth.NewCSRFToken(sessionToken)
	if err != nil {
		t.Fatal(err)
	}
	h := requireCSRF(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, tc := range []struct {
		name, csrf string
		want int
	}{
		{"valid", csrfToken, 204},
		{"missing", "", 403},
		{"invalid", "bad-token", 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
			req.AddCookie(&http.Cookie{Name: session.CookieName, Value: sessionToken})
			req.Header.Set("X-CSRF-Token", tc.csrf)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status %d, want %d", rec.Code, tc.want)
			}
		})
	}
}
