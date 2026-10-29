package session

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
)

// ConsumeChallenge atomically validates and consumes an issued login challenge.
// Issuance and email delivery are deliberately separate and not exposed over HTTP.
func (r Repository) ConsumeChallenge(ctx context.Context, token string) (string, error) {
 if r.DB == nil { return "", errors.New("authentication database unavailable") }
 digest, err := Digest(token)
 if err != nil { return "", ErrNotFound }
 var userID string
 err = r.DB.QueryRowContext(ctx, `
  UPDATE login_challenges
  SET consumed_at = clock_timestamp()
  WHERE token_digest = $1
    AND consumed_at IS NULL
    AND expires_at > clock_timestamp()
    AND attempts < 5
  RETURNING user_id`, digest[:]).Scan(&userID)
 if errors.Is(err, sql.ErrNoRows) { return "", ErrNotFound }
 if err != nil { return "", fmt.Errorf("consume login challenge: %w", err) }
 return userID, nil
}
