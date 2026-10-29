package main

import (
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/formschema"
 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func TestPublicFormPageEscapesUntrustedContent(t *testing.T) {
 form:=store.PublicForm{Title:"<script>alert(1)</script>",Schema:formschema.Schema{Fields:[]formschema.Field{{ID:"name",Type:"text",Label:"<img src=x onerror=alert(1)>",Required:true}}}}
 var b strings.Builder
 if err:=publicFormPage.Execute(&b,form);err!=nil {t.Fatal(err)}
 html:=b.String()
 if strings.Contains(html,"<script>alert(1)</script>") || strings.Contains(html,"<img src=x") {t.Fatal("unescaped content")}
 if !strings.Contains(html,"&lt;script&gt;") || !strings.Contains(html,"required") {t.Fatal("expected escaped content and required field")}
 if !strings.Contains(html,"Send melding") || !strings.Contains(html,"name=\"_website\"") {t.Fatal("submit and honeypot must be present")}
}

func TestPublicFormPageUnavailableWithoutDatabase(t *testing.T) {
 w:=httptest.NewRecorder()
 r:=httptest.NewRequest(http.MethodGet,"/f/00000000-0000-0000-0000-000000000001/demo/contact",nil)
 handler().ServeHTTP(w,r)
 if w.Code!=http.StatusNotFound {t.Fatalf("got %d want 404",w.Code)}
}
