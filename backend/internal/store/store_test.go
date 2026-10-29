package store

import "testing"

func TestValidateSlug(t *testing.T) {
	for _, v := range []string{"it", "my-project", "a123"} {
		if err := ValidateSlug(v); err != nil {
			t.Errorf("%q: %v", v, err)
		}
	}
	for _, v := range []string{"", "a", "UPPER", "bad_slug", "bad slug"} {
		if err := ValidateSlug(v); err == nil {
			t.Errorf("expected error for %q", v)
		}
	}
}

func TestValidUUID(t *testing.T) {
	if !validUUID("123e4567-e89b-12d3-a456-426614174000") {
		t.Fatal("valid uuid rejected")
	}
	if validUUID("not-a-uuid") {
		t.Fatal("invalid uuid accepted")
	}
}
