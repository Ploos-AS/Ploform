package authz

import (
	"context"
	"errors"
	"testing"

	"github.com/Ploos-AS/Ploform/backend/internal/session"
)

type fakeSessions struct {
	identity session.Identity
	err      error
}

func (f fakeSessions) Lookup(context.Context, string) (session.Identity, error) {
	return f.identity, f.err
}

type fakeMemberships struct {
	role Role
	err  error
}

func (f fakeMemberships) Resolve(context.Context, string, string) (Role, error) {
	return f.role, f.err
}

func TestAuthorize(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name    string
		session fakeSessions
		member  fakeMemberships
		perm    Permission
		wantErr error
	}{
		{"owner", fakeSessions{identity: session.Identity{UserID: "user-1"}}, fakeMemberships{role: Owner}, ManageOrg, nil},
		{"viewer denied", fakeSessions{identity: session.Identity{UserID: "user-1"}}, fakeMemberships{role: Viewer}, ManageForms, ErrForbidden},
		{"missing membership", fakeSessions{identity: session.Identity{UserID: "user-1"}}, fakeMemberships{err: ErrNoMembership}, ReadResponses, ErrForbidden},
		{"invalid token", fakeSessions{err: session.ErrNotFound}, fakeMemberships{role: Owner}, ManageOrg, ErrUnauthenticated},
		{"empty identity", fakeSessions{}, fakeMemberships{role: Owner}, ManageOrg, ErrUnauthenticated},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := Authorizer{Sessions: tc.session, Memberships: tc.member}
			principal, err := a.Authorize(ctx, "test-token", "org-1", tc.perm)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Authorize error = %v, want %v", err, tc.wantErr)
			}
			if err == nil && (principal.UserID != "user-1" || principal.OrganizationID != "org-1") {
				t.Fatalf("unexpected principal: %+v", principal)
			}
		})
	}
}

func TestAuthorizeFailsClosedWithoutDependencies(t *testing.T) {
	_, err := (Authorizer{}).Authorize(context.Background(), "token", "org", ManageForms)
	if !errors.Is(err, ErrUnauthenticated) {
		t.Fatalf("expected unauthenticated, got %v", err)
	}
}
