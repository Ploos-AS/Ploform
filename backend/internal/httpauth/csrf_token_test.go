package httpauth

import (
	"strings"
	"testing"
)

func TestCSRFTokenBoundToSession(t *testing.T) {
	token, err := NewCSRFToken("session-a")
	if err != nil {
		t.Fatal(err)
	}
	if !ValidCSRFToken("session-a", token) {
		t.Fatal("valid token rejected")
	}
	if ValidCSRFToken("session-b", token) {
		t.Fatal("token accepted for another session")
	}
	if ValidCSRFToken("session-a", token+"x") || ValidCSRFToken("", token) {
		t.Fatal("malformed or unbound token accepted")
	}
	if _, err := NewCSRFToken(""); err == nil {
		t.Fatal("accepted missing session")
	}
}

func TestCSRFTokensHaveFreshNonces(t *testing.T) {
	a, err := NewCSRFToken("session")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewCSRFToken("session")
	if err != nil {
		t.Fatal(err)
	}
	if a == b || strings.TrimSpace(a) == "" {
		t.Fatal("CSRF tokens not randomized")
	}
}
