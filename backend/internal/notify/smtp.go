// Package notify provides opt-in SMTP delivery for form notifications.
// Recipient addresses and credentials must be configured server-side.
package notify

import (
 "context"
 "crypto/tls"
 "errors"
 "fmt"
 "net"
 "net/mail"
 "mime"
 "net/smtp"
 "strconv"
 "strings"
 "time"
)

type SMTPConfig struct {
 Host string
 Port int
 Username string
 Password string
 From string
 To string
}

func (c SMTPConfig) Validate() error {
 if c.Host=="" || strings.ContainsAny(c.Host,"\r\n") {return errors.New("SMTP host required")}
 if c.Port<1 || c.Port>65535 {return errors.New("invalid SMTP port")}
 if c.Username=="" || c.Password=="" {return errors.New("SMTP credentials required")}
 from,err:=mail.ParseAddress(c.From);if err!=nil || from.Address!=c.From {return errors.New("invalid SMTP sender")}
 to,err:=mail.ParseAddress(c.To);if err!=nil || to.Address!=c.To {return errors.New("invalid SMTP recipient")}
 return nil
}

// Send uses mandatory STARTTLS and authenticated SMTP. It refuses plaintext fallback.
// This method is intentionally synchronous; production should enqueue notifications
// after persisting a submission, with retries and idempotent delivery.
func (c SMTPConfig) Send(ctx context.Context, subject, body string) error {
 if err:=c.Validate();err!=nil{return err}
 if strings.ContainsAny(subject,"\r\n") {return errors.New("invalid subject")}
 if ctx==nil {return errors.New("context required")}
 addr:=net.JoinHostPort(c.Host,strconv.Itoa(c.Port))
 dialer:=net.Dialer{Timeout:10*time.Second}
 conn,err:=dialer.DialContext(ctx,"tcp",addr)
 if err!=nil{return fmt.Errorf("SMTP connection failed: %w",err)}
 defer conn.Close()
 deadline:=time.Now().Add(20*time.Second)
 if d,ok:=ctx.Deadline();ok && d.Before(deadline){deadline=d}
 if err:=conn.SetDeadline(deadline);err!=nil{return err}
 client,err:=smtp.NewClient(conn,c.Host)
 if err!=nil{return fmt.Errorf("SMTP handshake failed: %w",err)}
 defer client.Close()
 ok,_:=client.Extension("STARTTLS")
 if !ok{return errors.New("SMTP server does not advertise STARTTLS")}
 if err:=client.StartTLS(&tls.Config{ServerName:c.Host,MinVersion:tls.VersionTLS12});err!=nil{return fmt.Errorf("SMTP TLS failed: %w",err)}
 if err:=client.Auth(smtp.PlainAuth("",c.Username,c.Password,c.Host));err!=nil{return fmt.Errorf("SMTP authentication failed: %w",err)}
 if err:=client.Mail(c.From);err!=nil{return fmt.Errorf("SMTP sender rejected: %w",err)}
 if err:=client.Rcpt(c.To);err!=nil{return fmt.Errorf("SMTP recipient rejected: %w",err)}
 writer,err:=client.Data();if err!=nil{return fmt.Errorf("SMTP data rejected: %w",err)}
 // Use fixed headers, with no user-controlled recipient or sender headers.
 payload:="From: "+c.From+"\r\nTo: "+c.To+"\r\nSubject: "+mime.QEncoding.Encode("utf-8",subject)+"\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n"+strings.ReplaceAll(strings.ReplaceAll(body,"\r\n","\n"),"\n","\r\n")
 _,writeErr:=writer.Write([]byte(payload))
 closeErr:=writer.Close()
 if writeErr!=nil{return fmt.Errorf("SMTP write failed: %w",writeErr)}
 if closeErr!=nil{return fmt.Errorf("SMTP message rejected: %w",closeErr)}
 return client.Quit()
}
