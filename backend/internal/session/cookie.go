package session

import (
	"errors"
	"net/http"
	"time"
)

const CookieName = "__Host-ploform_session"

// SetCookie issues a host-only session cookie. Call only after successful
// authentication over HTTPS. Session creation is handled separately.
func SetCookie(w http.ResponseWriter, token string, expiresAt time.Time) error {
	if _, err := Digest(token); err != nil {
		return errors.New("refusing invalid session cookie")
	}
	if !expiresAt.After(time.Now()) {
		return errors.New("refusing expired session cookie")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

// ClearCookie removes the browser cookie. The server must revoke the session
// digest in the database separately; deleting the cookie is not logout alone.
func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
