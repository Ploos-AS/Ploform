package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrNotFound = errors.New("session not found")

// Repository uses a dedicated, narrowly privileged authentication connection.
// Never connect it using the migration owner or a superuser.
type Repository struct {
	DB *sql.DB
}

type Identity struct {
	UserID string
	ExpiresAt time.Time
}

// Create persists only the token digest. The caller must have authenticated
// userID already, and must deliver the raw token through a secure channel.
func (r Repository) Create(ctx context.Context, userID string, expiresAt time.Time) (string, error) {
	if r.DB == nil {
		return "", errors.New("authentication database unavailable")
	}
	if userID == "" || !expiresAt.After(time.Now()) {
		return "", errors.New("invalid session parameters")
	}
	token, digest, err := NewToken()
	if err != nil {
		return "", err
	}
	_, err = r.DB.ExecContext(ctx,
		"INSERT INTO sessions (user_id, token_digest, expires_at) VALUES ($1, $2, $3)",
		userID, digest[:], expiresAt)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	return token, nil
}

// Lookup checks expiry and revocation in the database on every request.
// No result is returned for revoked or expired sessions.
func (r Repository) Lookup(ctx context.Context, token string) (Identity, error) {
	var identity Identity
	if r.DB == nil {
		return identity, errors.New("authentication database unavailable")
	}
	digest, err := Digest(token)
	if err != nil {
		return identity, ErrNotFound
	}
	err = r.DB.QueryRowContext(ctx,
		`SELECT user_id, expires_at FROM sessions
		  WHERE token_digest = $1 AND revoked_at IS NULL AND expires_at > now()`,
		digest[:]).Scan(&identity.UserID, &identity.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Identity{}, ErrNotFound
	}
	if err != nil {
		return Identity{}, fmt.Errorf("lookup session: %w", err)
	}
	return identity, nil
}

// Revoke invalidates the current token. Unknown tokens are treated as absent.
func (r Repository) Revoke(ctx context.Context, token string) error {
	if r.DB == nil {
		return errors.New("authentication database unavailable")
	}
	digest, err := Digest(token)
	if err != nil {
		return ErrNotFound
	}
	result, err := r.DB.ExecContext(ctx,
		"UPDATE sessions SET revoked_at = now() WHERE token_digest = $1 AND revoked_at IS NULL",
		digest[:])
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}
