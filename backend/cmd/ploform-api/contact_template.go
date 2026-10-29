package main

import (
 "encoding/json"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

// contactTemplateHandler exposes the built-in preset for the future
// authenticated create-from-template flow. It does not create a form.
func contactTemplateHandler(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "application/json")
 w.Header().Set("Cache-Control", "no-store")
 _ = json.NewEncoder(w).Encode(map[string]any{
  "id": "contact",
  "name": "Kontaktskjema",
  "schema": formschema.Contact(),
 })
}
