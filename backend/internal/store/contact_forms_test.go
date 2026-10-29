package store

import (
 "context"
 "errors"
 "testing"
)

func TestCreateContactFormRejectsInvalidInputBeforeDatabase(t *testing.T) {
 s := Store{}
 cases := []struct {
  name, workspace, slug, title string
 }{
  {"invalid workspace", "not-a-uuid", "contact", "Kontakt"},
  {"invalid slug", "00000000-0000-0000-0000-000000000002", "Bad Slug", "Kontakt"},
  {"empty title", "00000000-0000-0000-0000-000000000002", "contact", "  "},
  {"long title", "00000000-0000-0000-0000-000000000002", "contact", string(make([]byte, 201))},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   _, _, err := s.CreateContactForm(context.Background(), "00000000-0000-0000-0000-000000000001", tc.workspace, tc.slug, tc.title)
   if !errors.Is(err, ErrInvalidForm) { t.Fatalf("got %v, want ErrInvalidForm", err) }
  })
 }
}
