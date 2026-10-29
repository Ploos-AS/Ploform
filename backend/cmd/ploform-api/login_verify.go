package main

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "time"

 "github.com/Ploos-AS/Ploform/backend/internal/session"
)

type challengeRedeemer interface {
 RedeemChallenge(context.Context,string,time.Duration)(string,time.Time,error)
}

// verifyLoginHandler accepts only a JSON POST. Route mounting is deferred
// until rate limiting and same-origin protection are configured.
func verifyLoginHandler(redeemer challengeRedeemer) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  if r.Method!=http.MethodPost {w.Header().Set("Allow","POST");http.Error(w,"method not allowed",http.StatusMethodNotAllowed);return}
  if r.Header.Get("Content-Type")!="application/json" {http.Error(w,"JSON required",http.StatusUnsupportedMediaType);return}
  r.Body=http.MaxBytesReader(w,r.Body,4096)
  var body struct{Token string `json:"token"`}
  if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil || body.Token=="" {http.Error(w,"invalid request",http.StatusBadRequest);return}
  token,expiry,err:=redeemer.RedeemChallenge(r.Context(),body.Token,24*time.Hour)
  if err!=nil {
   if errors.Is(err,session.ErrNotFound){http.Error(w,"invalid or expired link",http.StatusUnauthorized);return}
   http.Error(w,"login unavailable",http.StatusServiceUnavailable);return
  }
  if err:=session.SetCookie(w,token,expiry);err!=nil{http.Error(w,"login unavailable",http.StatusServiceUnavailable);return}
  w.WriteHeader(http.StatusNoContent)
 })
}
