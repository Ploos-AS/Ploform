package authz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoMembership = errors.New("organization membership not found")

// MembershipRepository must use an authorized, narrowly privileged identity
// connection. It is not accessible to the public tenant application role.
type MembershipRepository struct {
	DB *sql.DB
}

// Resolve verifies membership from server-side storage. orgID may originate
// from a route, but the role must never originate from client input.
func (r MembershipRepository) Resolve(ctx context.Context, userID, orgID string) (Role, error) {
	if r.DB == nil {
		return "", errors.New("membership database unavailable")
	}
	if userID == "" || orgID == "" {
		return "", ErrNoMembership
	}
	var role string
	err := r.DB.QueryRowContext(ctx,
		"SELECT role FROM organization_memberships WHERE user_id = $1 AND organization_id = $2",
		userID, orgID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoMembership
	}
	if err != nil {
		return "", fmt.Errorf("resolve membership: %w", err)
	}
	switch Role(role) {
	case Owner, Admin, Editor, Viewer:
		return Role(role), nil
	default:
		return "", ErrNoMembership
	}
}
