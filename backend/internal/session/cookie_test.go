package session

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCookieSecurityFlags(t *testing.T) {
	token, _, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	if err := SetCookie(rec, token, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	c := cookies[0]
	if c.Name != CookieName || c.Path != "/" || !c.Secure || !c.HttpOnly || c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("insecure cookie: %+v", c)
	}
	if !strings.Contains(rec.Header().Get("Set-Cookie"), "__Host-") {
		t.Fatal("missing host-only cookie prefix")
	}
}

func TestRejectInvalidCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	if err := SetCookie(rec, "invalid", time.Now().Add(time.Hour)); err == nil {
		t.Fatal("accepted invalid cookie")
	}
	if rec.Header().Get("Set-Cookie") != "" {
		t.Fatal("set cookie after error")
	}
}

func TestClearCookie(t *testing.T) {
	rec := httptest.NewRecorder()
	ClearCookie(rec)
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].MaxAge != -1 {
		t.Fatalf("cookie not cleared: %+v", cookies)
	}
}
