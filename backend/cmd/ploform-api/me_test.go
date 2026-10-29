package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMeRouteDisabledWithoutAuthDatabase(t *testing.T) {
	path := "/api/me"
	for _, handler := range []http.Handler{handlerWithDatabases(nil, nil)} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status %d, want 404", rec.Code)
		}
	}
}
