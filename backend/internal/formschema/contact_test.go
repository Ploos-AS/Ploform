package formschema

import "testing"

func TestContactTemplate(t *testing.T) {
 s := Contact()
 if err := Validate(s); err != nil {
  t.Fatalf("contact template must be valid: %v", err)
 }
 want := []struct {
  id, kind string
  required bool
 }{
  {"name", "text", true},
  {"email", "email", true},
  {"subject", "text", false},
  {"message", "textarea", true},
 }
 if len(s.Fields) != len(want) {
  t.Fatalf("got %d fields, want %d", len(s.Fields), len(want))
 }
 for i, field := range s.Fields {
  if field.ID != want[i].id || field.Type != want[i].kind || field.Required != want[i].required {
   t.Errorf("field %d: got %+v, want %+v", i, field, want[i])
  }
 }
}

func TestContactTemplateIndependentInstances(t *testing.T) {
 first := Contact()
 first.Fields[0].Label = "Modified"
 second := Contact()
 if second.Fields[0].Label != "Navn" {
  t.Fatal("contact templates must not share mutable state")
 }
}
