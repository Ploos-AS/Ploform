package authz

import (
	"context"
	"errors"

	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

// SessionLookup and MembershipLookup allow authentication and authorization
// to be tested without a database or HTTP server.
type SessionLookup interface {
	Lookup(ctx context.Context, token string) (session.Identity, error)
}

type MembershipLookup interface {
	Resolve(ctx context.Context, userID, orgID string) (Role, error)
}

type Principal struct {
	UserID         string
	OrganizationID string
	Role           Role
}

var ErrUnauthenticated = errors.New("authentication required")

// Authorizer never trusts a role supplied by the caller. The organization ID
// is only a requested scope; the membership lookup must confirm it.
type Authorizer struct {
	Sessions    SessionLookup
	Memberships MembershipLookup
}

func (a Authorizer) Authorize(ctx context.Context, token, organizationID string, permission Permission) (Principal, error) {
	if a.Sessions == nil || a.Memberships == nil || token == "" {
		return Principal{}, ErrUnauthenticated
	}
	identity, err := a.Sessions.Lookup(ctx, token)
	if err != nil || identity.UserID == "" {
		return Principal{}, ErrUnauthenticated
	}
	role, err := a.Memberships.Resolve(ctx, identity.UserID, organizationID)
	if err != nil {
		return Principal{}, ErrForbidden
	}
	if err := Require(role, permission); err != nil {
		return Principal{}, err
	}
	return Principal{UserID: identity.UserID, OrganizationID: organizationID, Role: role}, nil
}
