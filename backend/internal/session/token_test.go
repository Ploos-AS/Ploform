package session

import (
	"testing"
	"time"
)

func TestNewTokenAndDigest(t *testing.T) {
	a, ah, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}
	b, bh, err := NewToken()
	if err != nil {
		t.Fatal(err)
	}
	if a == b || ah == bh {
		t.Fatal("tokens must be unique")
	}
	got, err := Digest(a)
	if err != nil {
		t.Fatal(err)
	}
	if got != ah {
		t.Fatal("digest mismatch")
	}
}

func TestInvalidTokens(t *testing.T) {
	for _, token := range []string{"", "not-a-token", "AA", "aGVsbG8="} {
		if _, err := Digest(token); err == nil {
			t.Errorf("accepted %q", token)
		}
	}
}

func TestActive(t *testing.T) {
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	if err := Active(now, now.Add(time.Minute), false); err != nil {
		t.Fatal(err)
	}
	if err := Active(now, now, false); err == nil {
		t.Fatal("accepted expired session")
	}
	if err := Active(now, now.Add(time.Minute), true); err == nil {
		t.Fatal("accepted revoked session")
	}
}
