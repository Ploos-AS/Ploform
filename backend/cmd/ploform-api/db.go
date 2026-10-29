package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// connectDatabase verifies connectivity before returning a pooled connection.
// Only the restricted application role must be used for API traffic.
func connectDatabase(ctx context.Context) (*sql.DB, error) {
	dsn, err := databaseDSN("PLOFORM_DATABASE_URL", "PLOFORM_APP_PASSWORD_FILE", "ploform_app")
	if err != nil { return nil, err }
	if dsn == "" {
		return nil, errors.New("PLOFORM_DATABASE_URL is required for database mode")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	var currentUser string
	if err := db.QueryRowContext(ctx, "SELECT current_user").Scan(&currentUser); err != nil || currentUser != "ploform_app" {
		_ = db.Close()
		return nil, errors.New("application database must connect as ploform_app")
	}
	return db, nil
}
