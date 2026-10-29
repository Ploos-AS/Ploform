package session

import (
 "context"
 "errors"
 "net/url"
 "strings"
)

// ChallengeCreator issues an opaque, short-lived credential for an already
// authenticated/verified user ID.
type ChallengeCreator interface {
 IssueChallenge(context.Context,string)(string,error)
}

// LoginDelivery coordinates issuance and delivery. The caller must resolve
// the user from a trusted source and enforce request throttling before calling.
type LoginDelivery struct {
 Issuer ChallengeCreator
 Mailer MailSender
 PublicOrigin string
}

func (d LoginDelivery) Send(ctx context.Context,userID,email string) error {
 if d.Issuer==nil || d.Mailer==nil {return errors.New("login delivery unavailable")}
 origin,err:=url.Parse(d.PublicOrigin)
 if err!=nil || origin.Scheme!="https" || origin.Host=="" || origin.User!=nil || origin.RawQuery!="" || origin.Fragment!="" || origin.Path!="" {
  return errors.New("login origin must be an HTTPS origin")
 }
 if strings.ContainsAny(email,"\r\n") || !strings.Contains(email,"@") {
  return errors.New("invalid recipient")
 }
 token,err:=d.Issuer.IssueChallenge(ctx,userID)
 if err!=nil{return err}
 // A fragment keeps the bearer token out of HTTP request lines and proxy logs.
 // The future verification page must read it client-side and submit it via POST.
 link:=origin.String()+"/login/verify#token="+url.QueryEscape(token)
 return d.Mailer.SendLoginLink(ctx,email,link)
}
