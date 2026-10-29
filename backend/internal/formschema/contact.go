package formschema

// Contact returns a fresh, editable schema for a simple contact form.
// It uses only the general form schema's supported field types.
func Contact() Schema {
 return Schema{Fields: []Field{
  {ID: "name", Type: "text", Label: "Navn", Required: true},
  {ID: "email", Type: "email", Label: "E-post", Required: true},
  {ID: "subject", Type: "text", Label: "Emne"},
  {ID: "message", Type: "textarea", Label: "Melding", Required: true},
 }}
}
