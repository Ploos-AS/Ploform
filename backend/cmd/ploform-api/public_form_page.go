package main

import (
 "errors"
 "html/template"
 "net/http"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

var publicFormPage = template.Must(template.New("form").Parse(`<!doctype html>
<html lang="no"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="robots" content="noindex"><title>{{.Title}}</title>
<style>
:root{color-scheme:light dark;font-family:system-ui,sans-serif}body{max-width:42rem;margin:3rem auto;padding:0 1rem;line-height:1.5}
label{display:block;font-weight:600;margin:1.2rem 0 .3rem}
input,textarea,select{box-sizing:border-box;width:100%;font:inherit;padding:.65rem;border:1px solid #888;border-radius:.35rem}
textarea{min-height:9rem}fieldset{border:0;padding:0;margin:1rem 0}
button{font:inherit;padding:.7rem 1.3rem;cursor:pointer}
small{display:block;margin-top:.5rem}
</style></head><body><main><h1>{{.Title}}</h1>
<form method="post" action="" accept-charset="utf-8">
{{range .Schema.Fields}}
{{if eq .Type "checkbox"}}<label><input style="width:auto" type="checkbox" name="{{.ID}}" {{if .Required}}required{{end}}> {{.Label}}</label>
{{else}}
<label for="field-{{.ID}}">{{.Label}}{{if .Required}} <span aria-label="obligatorisk">*</span>{{end}}</label>
{{if eq .Type "textarea"}}<textarea id="field-{{.ID}}" name="{{.ID}}" {{if .Required}}required{{end}}></textarea>
{{else if eq .Type "select"}}<select id="field-{{.ID}}" name="{{.ID}}" {{if .Required}}required{{end}}><option value="">Velg …</option>{{range .Options}}<option value="{{.}}">{{.}}</option>{{end}}</select>
{{else if eq .Type "multiselect"}}<select id="field-{{.ID}}" name="{{.ID}}" multiple {{if .Required}}required{{end}}>{{range .Options}}<option value="{{.}}">{{.}}</option>{{end}}</select>
{{else}}<input id="field-{{.ID}}" name="{{.ID}}" type="{{.Type}}" {{if .Required}}required{{end}}>{{end}}
{{end}}{{end}}
<p><small>Felter merket * er obligatoriske.</small></p>
<button type="submit" disabled aria-disabled="true">Send melding (kommer snart)</button>
<p role="status">Innsending er ikke aktivert ennå.</p>
</form></main></body></html>`))

func publicFormPageHandler(s store.Store) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  form,err:=s.GetPublicForm(r.Context(),r.PathValue("organizationID"),r.PathValue("workspaceSlug"),r.PathValue("formSlug"))
  if err!=nil {
   if errors.Is(err,store.ErrFormNotFound) {http.NotFound(w,r);return}
   http.Error(w,"form unavailable",http.StatusServiceUnavailable);return
  }
  // All schema field types are checked by GetPublicForm before rendering.
  // Keep the type reference here to make the shared schema contract explicit.
  _ = formschema.Schema{}
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Content-Security-Policy","default-src 'none'; style-src 'unsafe-inline'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'")
  if err:=publicFormPage.Execute(w,form);err!=nil {return}
 })
}
