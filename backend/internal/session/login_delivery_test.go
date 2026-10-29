package session

import (
 "context"
 "errors"
 "strings"
 "testing"
)

type stubIssuer struct{token string;err error;calls int}
func (s *stubIssuer) IssueChallenge(_ context.Context,_ string)(string,error){s.calls++;return s.token,s.err}
type stubMailer struct{recipient,link string;err error;calls int}
func (m *stubMailer) SendLoginLink(_ context.Context,recipient,link string)error{m.calls++;m.recipient=recipient;m.link=link;return m.err}

func TestLoginDelivery(t *testing.T){
 issuer:=&stubIssuer{token:"secret-test-token"}
 mailer:=&stubMailer{}
 d:=LoginDelivery{Issuer:issuer,Mailer:mailer,PublicOrigin:"https://forms.example.test"}
 if err:=d.Send(context.Background(),"user-id","person@example.test");err!=nil{t.Fatal(err)}
 if issuer.calls!=1||mailer.calls!=1||mailer.recipient!="person@example.test"{t.Fatal("expected one issue and one delivery")}
 if !strings.HasPrefix(mailer.link,"https://forms.example.test/login/verify#token="){t.Fatalf("unexpected link shape")}
 if strings.Contains(strings.Split(mailer.link,"#")[0],"secret-test-token"){t.Fatal("token leaked into request URL")}
}
func TestLoginDeliveryFailsClosed(t *testing.T){
 for _,origin:=range []string{"http://forms.example.test","https://forms.example.test/extra","https://user:pass@forms.example.test",""}{
  t.Run(origin,func(t *testing.T){
   issuer:=&stubIssuer{token:"token"};mailer:=&stubMailer{}
   err:= (LoginDelivery{Issuer:issuer,Mailer:mailer,PublicOrigin:origin}).Send(context.Background(),"user","person@example.test")
   if err==nil||issuer.calls!=0||mailer.calls!=0{t.Fatalf("unsafe origin accepted: %q",origin)}
  })
 }
 issuer:=&stubIssuer{err:errors.New("database unavailable")};mailer:=&stubMailer{}
 err:=(LoginDelivery{Issuer:issuer,Mailer:mailer,PublicOrigin:"https://forms.example.test"}).Send(context.Background(),"user","person@example.test")
 if err==nil||mailer.calls!=0{t.Fatal("mailer must not run after issuance failure")}
}
