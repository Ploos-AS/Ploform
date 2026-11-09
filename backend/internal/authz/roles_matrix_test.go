package authz

import (
    "errors"
    "testing"
)

func TestAllowedRolePermissionMatrix(t *testing.T) {
    permissions := []Permission{ManageMembers, ManageForms, ReadResponses, ManageOrg}
    cases := []struct {
        role Role
        allowed map[Permission]bool
    }{
        {Owner, map[Permission]bool{ManageMembers:true, ManageForms:true, ReadResponses:true, ManageOrg:true}},
        {Admin, map[Permission]bool{ManageMembers:true, ManageForms:true, ReadResponses:true}},
        {Editor, map[Permission]bool{ManageForms:true, ReadResponses:true}},
        {Viewer, map[Permission]bool{ReadResponses:true}},
        {Role("unknown"), map[Permission]bool{}},
        {Role(""), map[Permission]bool{}},
    }
    for _, tc := range cases {
        for _, permission := range permissions {
            t.Run(string(tc.role)+"/"+string(permission), func(t *testing.T) {
                want := tc.allowed[permission]
                if got := Allowed(tc.role, permission); got != want {
                    t.Fatalf("Allowed(%q,%q)=%t; want %t", tc.role, permission, got, want)
                }
                err := Require(tc.role, permission)
                if want && err != nil {
                    t.Fatalf("Require(%q,%q) unexpected error: %v", tc.role, permission, err)
                }
                if !want && !errors.Is(err, ErrForbidden) {
                    t.Fatalf("Require(%q,%q) error=%v; want ErrForbidden", tc.role, permission, err)
                }
            })
        }
        if Allowed(tc.role, Permission("unknown")) {
            t.Errorf("role %q allowed unknown permission", tc.role)
        }
    }
}
