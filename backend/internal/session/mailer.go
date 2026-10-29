package session

import (
 "context"
 "errors"
 "fmt"
 "net"
 "net/smtp"
 "strings"
 "time"
)

// MailSender is deliberately independent of the HTTP login flow.
type MailSender interface {
 SendLoginLink(ctx context.Context, recipient, link string) error
}

// SMTPMailer delivers one-time links using a configured authenticated SMTP relay.
// Use TLS on the connection; never send bearer credentials over plaintext SMTP.
type SMTPMailer struct {
 Address string // host:port
 Username string
 Password string
 From string
}

func (m SMTPMailer) SendLoginLink(ctx context.Context, recipient, link string) error {
 if m.Address=="" || m.Username=="" || m.Password=="" || m.From=="" {
  return errors.New("SMTP mailer is not configured")
 }
 if strings.ContainsAny(recipient,"\r\n") || strings.ContainsAny(m.From,"\r\n") || strings.ContainsAny(link,"\r\n") {
  return errors.New("invalid mail header")
 }
 host,_,err:=net.SplitHostPort(m.Address)
 if err!=nil{return fmt.Errorf("invalid SMTP address: %w",err)}
 if host==""{return errors.New("SMTP host required")}
 // This transport requires implicit TLS (typically port 465).
 // STARTTLS support can be added separately with explicit downgrade protection.
 dialer:=net.Dialer{Timeout:10*time.Second}
 conn,err:=dialer.DialContext(ctx,"tcp",m.Address)
 if err!=nil{return fmt.Errorf("SMTP connection: %w",err)}
 defer conn.Close()
 // smtp.SendMail does not provide an implicit TLS connection and must not
 // be used here. Fail closed until a TLS transport is implemented.
 return errors.New("implicit TLS SMTP transport not implemented")
}
