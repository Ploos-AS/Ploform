package formschema

import (
 "errors"
 "regexp"
)

var ErrInvalid = errors.New("invalid form schema")
var fieldID = regexp.MustCompile("^[a-z][a-z0-9_]{0,63}$")

func Validate(s Schema) error {
 if len(s.Fields) == 0 || len(s.Fields) > 100 { return ErrInvalid }
 seen := map[string]bool{}
 for _, f := range s.Fields {
  if !fieldID.MatchString(f.ID) || seen[f.ID] || len(f.Label) == 0 || len(f.Label) > 200 { return ErrInvalid }
  seen[f.ID] = true
  switch f.Type {
  case "text", "textarea", "email", "number", "date", "checkbox":
   if len(f.Options) != 0 { return ErrInvalid }
  case "select", "multiselect":
   if len(f.Options) < 2 || len(f.Options) > 50 { return ErrInvalid }
   options := map[string]bool{}
   for _, option := range f.Options {
    if len(option) == 0 || len(option) > 200 || options[option] { return ErrInvalid }
    options[option] = true
   }
  default:
   return ErrInvalid
  }
 }
 return nil
}
