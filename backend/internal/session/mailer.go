package session

import (
 "context"
 "crypto/tls"
 "errors"
 "fmt"
 "net"
 "net/smtp"
 "strings"
 "time"
)

type MailSender interface {
 SendLoginLink(ctx context.Context, recipient, link string) error
}

// SMTPMailer requires implicit TLS (typically port 465). No plaintext fallback.
type SMTPMailer struct {
 Address string
 Username string
 Password string
 From string
}

func (m SMTPMailer) SendLoginLink(ctx context.Context, recipient, link string) error {
 if m.Address=="" || m.Username=="" || m.Password=="" || m.From=="" {
  return errors.New("SMTP mailer is not configured")
 }
 if strings.ContainsAny(recipient,"\r\n") || strings.ContainsAny(m.From,"\r\n") || strings.ContainsAny(link,"\r\n") {
  return errors.New("invalid mail input")
 }
 host,_,err:=net.SplitHostPort(m.Address)
 if err!=nil || host=="" {return errors.New("invalid SMTP address")}
 dialer:=net.Dialer{Timeout:10*time.Second}
 conn,err:=tls.DialWithDialer(&dialer,"tcp",m.Address,&tls.Config{ServerName:host,MinVersion:tls.VersionTLS12})
 if err!=nil{return fmt.Errorf("secure SMTP connection: %w",err)}
 defer conn.Close()
 if err:=ctx.Err();err!=nil{return err}
 client,err:=smtp.NewClient(conn,host)
 if err!=nil{return fmt.Errorf("SMTP handshake: %w",err)}
 defer client.Close()
 if err:=client.Auth(smtp.PlainAuth("",m.Username,m.Password,host));err!=nil{return fmt.Errorf("SMTP authentication: %w",err)}
 if err:=client.Mail(m.From);err!=nil{return fmt.Errorf("SMTP sender: %w",err)}
 if err:=client.Rcpt(recipient);err!=nil{return fmt.Errorf("SMTP recipient: %w",err)}
 w,err:=client.Data()
 if err!=nil{return fmt.Errorf("SMTP DATA: %w",err)}
 body:="From: "+m.From+"\r\nTo: "+recipient+"\r\nSubject: Ploform sign-in link\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nUse this one-time link to sign in to Ploform:\r\n"+link+"\r\n\r\nThis link expires in 10 minutes. If you did not request it, ignore this email.\r\n"
 if _,err:=w.Write([]byte(body));err!=nil{w.Close();return fmt.Errorf("SMTP message: %w",err)}
 if err:=w.Close();err!=nil{return fmt.Errorf("SMTP message completion: %w",err)}
 return client.Quit()
}
