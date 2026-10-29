// Package submission validates untrusted public form values before persistence.
package submission

import (
 "errors"
 "net/mail"
 "strconv"
 "strings"
 "time"
 "unicode/utf8"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
)

var ErrInvalid = errors.New("invalid submission")

// Validate checks a decoded HTML form against the published schema.
// The returned map contains only recognized, normalized fields.
func Validate(schema formschema.Schema, values map[string][]string) (map[string]any, error) {
 if formschema.Validate(schema)!=nil { return nil,ErrInvalid }
 result:=make(map[string]any,len(schema.Fields))
 allowed:=make(map[string]bool,len(schema.Fields))
 for _,field:=range schema.Fields { allowed[field.ID]=true }
 for key:=range values { if !allowed[key] { return nil,ErrInvalid } }
 for _,f:=range schema.Fields {
  raw:=values[f.ID]
  if len(raw)>100 { return nil,ErrInvalid }
  if f.Type!="multiselect" && len(raw)>1 { return nil,ErrInvalid }
  if f.Type=="checkbox" {
   if len(raw)==0 {if f.Required{return nil,ErrInvalid};result[f.ID]=false;continue}
   if len(raw)!=1 || raw[0]!="on" {return nil,ErrInvalid}
   result[f.ID]=true
   continue
  }
  if f.Type=="multiselect" {
   if f.Required && len(raw)==0 {return nil,ErrInvalid}
   selected:=make([]string,0,len(raw))
   seen:=map[string]bool{}
   for _,v:=range raw {if !optionExists(f.Options,v)||seen[v]{return nil,ErrInvalid};seen[v]=true;selected=append(selected,v)}
   result[f.ID]=selected
   continue
  }
  if len(raw)==0 {if f.Required{return nil,ErrInvalid};continue}
  v:=strings.TrimSpace(raw[0])
  if !utf8.ValidString(v)||len(v)>10000 {return nil,ErrInvalid}
  if v=="" {if f.Required{return nil,ErrInvalid};continue}
  switch f.Type {
  case "text":
   if len(v)>500 {return nil,ErrInvalid}
  case "textarea":
   // Longer messages are supported, up to the global 10 KiB limit.
  case "email":
   if len(v)>254||strings.ContainsAny(v,"\r\n") {return nil,ErrInvalid}
   a,err:=mail.ParseAddress(v)
   if err!=nil||a.Address!=v||!strings.Contains(v,"@"){return nil,ErrInvalid}
  case "number":
   n,err:=strconv.ParseFloat(v,64)
   if err!=nil||strconv.FormatFloat(n,'g',-1,64)=="NaN"||n>1e308||n< -1e308 {return nil,ErrInvalid}
  case "date":
   d,err:=time.Parse("2006-01-02",v)
   if err!=nil||d.Format("2006-01-02")!=v{return nil,ErrInvalid}
  case "select":
   if !optionExists(f.Options,v){return nil,ErrInvalid}
  default:
   return nil,ErrInvalid
  }
  result[f.ID]=v
 }
 return result,nil
}
func optionExists(options []string,value string)bool{
 for _,option:=range options{if option==value{return true}}
 return false
}
