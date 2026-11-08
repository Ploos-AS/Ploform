package notify

import (
 "context"
 "strings"
 "testing"
)

func validConfig() SMTPConfig {
 return SMTPConfig{Host:"smtp.example.org",Port:587,Username:"user",Password:"secret",From:"sender@example.org",To:"recipient@example.org"}
}

func TestSMTPConfigValidation(t *testing.T) {
 if err:=validConfig().Validate();err!=nil{t.Fatal(err)}
 tests:=[]struct{name string;change func(*SMTPConfig)}{
  {"missing host",func(c *SMTPConfig){c.Host=""}},
  {"bad port",func(c *SMTPConfig){c.Port=0}},
  {"missing credentials",func(c *SMTPConfig){c.Password=""}},
  {"invalid sender",func(c *SMTPConfig){c.From="invalid"}},
  {"invalid recipient",func(c *SMTPConfig){c.To="bad\r\nBcc: attacker@example.org"}},
 }
 for _,tc:=range tests {t.Run(tc.name,func(t *testing.T){c:=validConfig();tc.change(&c);if err:=c.Validate();err==nil{t.Fatal("expected validation error")}})}
}

func TestSMTPRejectsHeaderInjectionBeforeConnecting(t *testing.T) {
 c:=validConfig()
 if err:=c.Send(context.Background(),"Hello\r\nBcc: attacker@example.org","body");err==nil || !strings.Contains(err.Error(),"subject"){t.Fatalf("expected subject rejection, got %v",err)}
}

func TestSMTPRequiresContext(t *testing.T) {
 if err:=validConfig().Send(nil,"Hello","body");err==nil{t.Fatal("expected context rejection")}
}
