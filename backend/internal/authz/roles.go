// Package authz defines permissions after identity and membership verification.
// It does not authenticate callers or resolve memberships.
package authz

import "errors"

type Role string
type Permission string

const (
	Owner  Role = "owner"
	Admin  Role = "admin"
	Editor Role = "editor"
	Viewer Role = "viewer"

	ManageMembers Permission = "members:manage"
	ManageForms   Permission = "forms:manage"
	ReadResponses Permission = "responses:read"
	ManageOrg     Permission = "organization:manage"
)

var ErrForbidden = errors.New("permission denied")

// Allowed checks permissions only for a role already verified by the server.
// Never accept a role value directly from a request as proof of authorization.
func Allowed(role Role, permission Permission) bool {
	switch role {
	case Owner:
		switch permission {
		case ManageMembers, ManageForms, ReadResponses, ManageOrg:
			return true
		}
	case Admin:
		switch permission {
		case ManageMembers, ManageForms, ReadResponses:
			return true
		}
	case Editor:
		return permission == ManageForms || permission == ReadResponses
	case Viewer:
		return permission == ReadResponses
	}
	return false
}

func Require(role Role, permission Permission) error {
	if !Allowed(role, permission) {
		return ErrForbidden
	}
	return nil
}
