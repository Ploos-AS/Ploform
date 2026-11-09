package notify

import (
 "context"
 "testing"
)

func TestWorkerFailsClosedWithoutDatabase(t *testing.T) {
 w:=Worker{SMTP:validConfig()}
 if err:=w.Run(context.Background());err==nil {t.Fatal("expected database requirement")}
}

func TestWorkerRejectsInvalidSMTPConfiguration(t *testing.T) {
 w:=Worker{}
 if err:=w.Run(context.Background());err==nil {t.Fatal("expected SMTP configuration error")}
}
