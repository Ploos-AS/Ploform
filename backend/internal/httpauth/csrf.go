package httpauth

import (
	"crypto/subtle"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
		if expected != "" && !validPublicOrigin(expected) {
			http.Error(w, "invalid server origin configuration", http.StatusServiceUnavailable)
			return
		}
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
		if err != nil || !validPublicOrigin(origin) {
			http.Error(w, "invalid origin", http.StatusForbidden)
			return
		}
		// Host must be validated against an operator-configured allowlist at
		// the reverse proxy before accepting public traffic.
		if expected != "" {
			if r.TLS == nil {
				if !TrustedProxyPeer(r.RemoteAddr, os.Getenv("PLOFORM_TRUSTED_PROXY_IPS")) {
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

// validPublicOrigin requires one canonical HTTPS origin without path, query,
// fragment, credentials or trailing slash. It rejects opaque URLs and ports
// that are not valid numeric TCP ports.
func validPublicOrigin(raw string) bool {
	if raw == "" || strings.TrimSpace(raw) != raw || strings.Contains(raw, ",") {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil ||
		u.Opaque != "" || u.Path != "" || u.RawQuery != "" ||
		u.Fragment != "" || u.RawFragment != "" {
		return false
	}
	if u.Hostname() == "" || strings.HasSuffix(raw, "/") {
		return false
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return false
		}
	}
	return true
}
