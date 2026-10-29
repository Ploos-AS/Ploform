package authz

import (
	"errors"
	"testing"
)

func TestPermissions(t *testing.T) {
	cases := []struct {
		role Role
		permission Permission
		want bool
	}{
		{Owner, ManageOrg, true},
		{Admin, ManageOrg, false},
		{Admin, ManageMembers, true},
		{Editor, ManageMembers, false},
		{Editor, ManageForms, true},
		{Viewer, ReadResponses, true},
		{Viewer, ManageForms, false},
		{Role("unknown"), ReadResponses, false},
		{Owner, Permission("unknown"), false},
	}
	for _, tc := range cases {
		if got := Allowed(tc.role, tc.permission); got != tc.want {
			t.Errorf("Allowed(%q, %q) = %v; want %v", tc.role, tc.permission, got, tc.want)
		}
	}
}

func TestRequireDenied(t *testing.T) {
	if err := Require(Viewer, ManageMembers); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
