package httpauth

import (
	"crypto/subtle"
	"net/http"
	"net"
	"net/url"
	"strings"
	"os"
)

// RequireSameOrigin rejects cross-origin unsafe requests. This is one
// defense for cookie-authenticated endpoints; it is not a complete CSRF
// solution for deployments with complex proxy configurations.
func RequireSameOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		expected := os.Getenv("PLOFORM_PUBLIC_ORIGIN")
		if expected == "" && r.TLS == nil {
			http.Error(w, "HTTPS required", http.StatusForbidden)
			return
		}
		origin := r.Header.Get("Origin")
		if origin == "" || strings.Contains(origin, ",") {
			http.Error(w, "origin required", http.StatusForbidden)
			return
		}
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
			http.Error(w, "invalid origin", http.StatusForbidden)
			return
		}
		// Host must be validated against an operator-configured allowlist at
		// the reverse proxy before accepting public traffic.
		if expected != "" {
			if r.TLS == nil {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil || net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() {
					http.Error(w, "trusted local proxy required", http.StatusForbidden)
					return
				}
			}
			if subtle.ConstantTimeCompare([]byte(origin), []byte(expected)) != 1 {
				http.Error(w, "cross-origin request denied", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		if subtle.ConstantTimeCompare([]byte(parsed.Host), []byte(r.Host)) != 1 {
			http.Error(w, "cross-origin request denied", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
