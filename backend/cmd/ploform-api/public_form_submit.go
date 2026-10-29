package main

import (
 "errors"
 "net/http"
 "sync"
 "time"

 "github.com/Ploos-AS/Ploform/backend/internal/store"
 "github.com/Ploos-AS/Ploform/backend/internal/submission"
)

// A conservative per-process limit. Production deployments with multiple
// replicas will need a shared limiter before accepting internet traffic.
type publicSubmitLimiter struct {
 mu sync.Mutex
 attempts map[string][]time.Time
}
func (l *publicSubmitLimiter) allow(key string) bool {
 l.mu.Lock()
 defer l.mu.Unlock()
 if l.attempts==nil {l.attempts=make(map[string][]time.Time)}
 now:=time.Now()
 // Periodic cleanup bounds memory from stale client addresses.
 if len(l.attempts)>10000 {
  for k,v:=range l.attempts {if len(v)==0||now.Sub(v[len(v)-1])>time.Minute {delete(l.attempts,k)}}
 }
 recent:=l.attempts[key][:0]
 for _,t:=range l.attempts[key] {if now.Sub(t)<time.Minute {recent=append(recent,t)}}
 if len(recent)>=5 {l.attempts[key]=recent;return false}
 l.attempts[key]=append(recent,now)
 return true
}

func publicFormSubmitHandler(s store.Store) http.Handler {
 limiter:=&publicSubmitLimiter{}
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("Content-Security-Policy","default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; frame-ancestors 'none'")
  if !limiter.allow(r.RemoteAddr) {http.Error(w,"For mange forsøk. Prøv igjen om litt.",http.StatusTooManyRequests);return}
  r.Body=http.MaxBytesReader(w,r.Body,32<<10)
  if err:=r.ParseForm();err!=nil {http.Error(w,"Ugyldig eller for stort skjema.",http.StatusBadRequest);return}
  if r.PostForm.Get("_website")!="" {http.Error(w,"Kunne ikke sende skjemaet.",http.StatusBadRequest);return}
  delete(r.PostForm,"_website")
  _,err:=s.SavePublicSubmission(r.Context(),r.PathValue("organizationID"),r.PathValue("workspaceSlug"),r.PathValue("formSlug"),r.PostForm)
  if err!=nil {
   if errors.Is(err,store.ErrFormNotFound){http.NotFound(w,r);return}
   if errors.Is(err,submission.ErrInvalid){http.Error(w,"Kontroller feltene og prøv igjen.",http.StatusBadRequest);return}
   http.Error(w,"Kunne ikke sende skjemaet. Prøv igjen senere.",http.StatusServiceUnavailable);return
  }
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  _,_=w.Write([]byte(`<!doctype html><html lang="no"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="robots" content="noindex"><title>Melding sendt</title></head><body><main><h1>Takk for meldingen!</h1><p>Henvendelsen er mottatt.</p></main></body></html>`))
 })
}
