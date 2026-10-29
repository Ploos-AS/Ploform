package formschema

import "testing"

func TestValidate(t *testing.T) {
 cases := []struct { name string; schema Schema; valid bool }{
  {"text", Schema{Fields: []Field{{ID:"name", Type:"text", Label:"Name"}}}, true},
  {"select", Schema{Fields: []Field{{ID:"choice", Type:"select", Label:"Choice", Options:[]string{"A","B"}}}}, true},
  {"empty", Schema{}, false},
  {"unknown type", Schema{Fields: []Field{{ID:"a", Type:"upload", Label:"File"}}}, false},
  {"duplicate ids", Schema{Fields: []Field{{ID:"a", Type:"text", Label:"One"},{ID:"a", Type:"text", Label:"Two"}}}, false},
  {"select missing options", Schema{Fields: []Field{{ID:"a", Type:"select", Label:"Choice"}}}, false},
  {"invalid id", Schema{Fields: []Field{{ID:"Bad ID", Type:"text", Label:"Name"}}}, false},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   err := Validate(tc.schema)
   if (err == nil) != tc.valid { t.Fatalf("Validate() error = %v, want valid=%v", err, tc.valid) }
  })
 }
}
