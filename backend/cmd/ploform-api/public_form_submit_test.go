package main

import (
 "net/http"
 "net/http/httptest"
 "net/url"
 "strings"
 "testing"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
)

func submitRequest(values url.Values) *http.Request {
 r:=httptest.NewRequest(http.MethodPost,"/f/00000000-0000-0000-0000-000000000001/demo/contact",strings.NewReader(values.Encode()))
 r.Header.Set("Content-Type","application/x-www-form-urlencoded")
 r.RemoteAddr="192.0.2.10:54321"
 r.SetPathValue("organizationID","00000000-0000-0000-0000-000000000001")
 r.SetPathValue("workspaceSlug","demo")
 r.SetPathValue("formSlug","contact")
 return r
}

func TestPublicSubmissionRejectsHoneypot(t *testing.T){
 h:=publicFormSubmitHandler(store.Store{})
 w:=httptest.NewRecorder()
 h.ServeHTTP(w,submitRequest(url.Values{"_website":{"https://spam.example"}}))
 if w.Code!=http.StatusBadRequest{t.Fatalf("status %d, want 400",w.Code)}
}

func TestPublicSubmissionRejectsOversizedBody(t *testing.T){
 h:=publicFormSubmitHandler(store.Store{})
 w:=httptest.NewRecorder()
 h.ServeHTTP(w,submitRequest(url.Values{"message":{strings.Repeat("x",33000)}}))
 if w.Code!=http.StatusBadRequest{t.Fatalf("status %d, want 400",w.Code)}
}

func TestPublicSubmissionRateLimitByIP(t *testing.T){
 h:=publicFormSubmitHandler(store.Store{})
 for i:=0;i<5;i++ {
  r:=submitRequest(url.Values{"_website":{"spam"}})
  r.RemoteAddr="192.0.2.10:"+string(rune('0'+i))+"1234"
  w:=httptest.NewRecorder()
  h.ServeHTTP(w,r)
  if w.Code!=http.StatusBadRequest{t.Fatalf("attempt %d: status %d",i,w.Code)}
 }
 w:=httptest.NewRecorder()
 h.ServeHTTP(w,submitRequest(url.Values{"_website":{"spam"}}))
 if w.Code!=http.StatusTooManyRequests{t.Fatalf("status %d, want 429",w.Code)}
}
