package httpauth

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSameOrigin(t *testing.T) {
	cases := []struct {
		name   string
		method string
		origin string
		secure bool
		want   int
	}{
		{"safe get", http.MethodGet, "", false, 204},
		{"valid post", http.MethodPost, "https://forms.example.test", true, 204},
		{"missing origin", http.MethodPost, "", true, 403},
		{"cross origin", http.MethodPost, "https://attacker.example.test", true, 403},
		{"insecure", http.MethodPost, "https://forms.example.test", false, 403},
		{"http origin", http.MethodPost, "http://forms.example.test", true, 403},
		{"origin with path", http.MethodPost, "https://forms.example.test/path", true, 403},
	}
	h := RequireSameOrigin(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "https://forms.example.test/api", nil)
			r.Host = "forms.example.test"
			if !tc.secure {
				r.TLS = nil
			} else {
				r.TLS = &tls.ConnectionState{}
			}
			if tc.origin != "" {
				r.Header.Set("Origin", tc.origin)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("status %d, want %d", w.Code, tc.want)
			}
		})
	}
}
