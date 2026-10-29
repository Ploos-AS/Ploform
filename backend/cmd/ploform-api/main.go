package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"status\":\"ok\"}\n"))
	})
	return mux
}

func main() {
	addr := os.Getenv("PLOFORM_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	// No administrative endpoints exist yet. Fail closed when a database
	// is configured but unavailable, rather than starting a degraded service.
	if os.Getenv("PLOFORM_DATABASE_URL") != "" {
		db, err := connectDatabase(context.Background())
		if err != nil {
			log.Fatalf("database startup check failed: %v", err)
		}
		defer db.Close()
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Ploform development API listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
