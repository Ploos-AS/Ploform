package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
)

// connectAuthDatabase must use ploform_auth, never the migration owner.
// The separate DSN prevents the tenant role from reading identity tables.
func connectAuthDatabase(ctx context.Context) (*sql.DB, error) {
	dsn := os.Getenv("PLOFORM_AUTH_DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("PLOFORM_AUTH_DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open authentication database: %w", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping authentication database: %w", err)
	}
	var currentUser string
	if err := db.QueryRowContext(ctx, "SELECT current_user").Scan(&currentUser); err != nil || currentUser != "ploform_auth" {
		_ = db.Close()
		return nil, errors.New("authentication database must connect as ploform_auth")
	}
	return db, nil
}
