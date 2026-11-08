package main

import (
 "testing"
 "time"
)

func TestLimiterAllowsSeparateClients(t *testing.T) {
 l:=&publicSubmitLimiter{}
 for i:=0;i<5;i++ {if !l.allow("192.0.2.10"){t.Fatalf("attempt %d rejected",i+1)}}
 if l.allow("192.0.2.10"){t.Fatal("sixth attempt should be rejected")}
 if !l.allow("192.0.2.11"){t.Fatal("independent client should be allowed")}
}

func TestLimiterExpiresOldAttempts(t *testing.T) {
 l:=&publicSubmitLimiter{attempts:map[string][]time.Time{"client":{time.Now().Add(-2*time.Minute)}}}
 if !l.allow("client"){t.Fatal("expired attempt should not block")}
 if got:=len(l.attempts["client"]);got!=1 {t.Fatalf("want one recent attempt, got %d",got)}
}
