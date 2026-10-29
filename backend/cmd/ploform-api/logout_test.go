package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type fakeRevoker struct {
	calls int
	err   error
}

func (f *fakeRevoker) Revoke(_ context.Context, _ string) error {
	f.calls++
	return f.err
}

func TestLogoutRevokesAndClearsCookie(t *testing.T) {
	revoker := &fakeRevoker{}
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "opaque-token"})
	rec := httptest.NewRecorder()
	logoutHandler(revoker).ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent || revoker.calls != 1 {
		t.Fatalf("status %d, revocations %d", rec.Code, revoker.calls)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != session.CookieName || cookies[0].MaxAge != -1 {
		t.Fatalf("session cookie not cleared: %+v", cookies)
	}
}

func TestLogoutDoesNotClearCookieWhenDatabaseFails(t *testing.T) {
	revoker := &fakeRevoker{err: errors.New("database unavailable")}
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	req.AddCookie(&http.Cookie{Name: session.CookieName, Value: "opaque-token"})
	rec := httptest.NewRecorder()
	logoutHandler(revoker).ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status %d, want 503", rec.Code)
	}
	if rec.Header().Get("Set-Cookie") != "" {
		t.Fatal("cookie cleared without confirmed revocation")
	}
}

func TestLogoutWithoutCookie(t *testing.T) {
	revoker := &fakeRevoker{}
	rec := httptest.NewRecorder()
	logoutHandler(revoker).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/logout", nil))
	if rec.Code != http.StatusUnauthorized || revoker.calls != 0 {
		t.Fatalf("status %d, revocations %d", rec.Code, revoker.calls)
	}
}
