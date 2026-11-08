package notify

import (
 "context"
 "errors"
 "log"
 "time"
)

// Worker polls the persistent outbox and sends due jobs using a configured
// SMTP transport. It must only run with server-side secrets.
type Worker struct {
 Outbox Outbox
 SMTP SMTPConfig
 PollInterval time.Duration
}

func (w Worker) Run(ctx context.Context) error {
 if ctx==nil {return errors.New("context required")}
 if err:=w.SMTP.Validate();err!=nil{return err}
 if w.Outbox.DB==nil {return errors.New("database required")}
 interval:=w.PollInterval
 if interval<=0 {interval=5*time.Second}
 ticker:=time.NewTicker(interval)
 defer ticker.Stop()
 for {
  if err:=ctx.Err();err!=nil{return nil}
  // Do not log email addresses, subjects, form values or SMTP credentials.
  err:=w.Outbox.ProcessOne(ctx,func(ctx context.Context,j Job)error{
   // Recipient must originate from trusted server-side configuration.
   // Per-form routing is intentionally not enabled yet.
   if j.Recipient!=w.SMTP.To {return errors.New("unapproved notification recipient")}
   return w.SMTP.Send(ctx,j.Subject,j.Body)
  })
  if err!=nil {
   if ctx.Err()!=nil{return nil}
   log.Print("notification worker cycle failed")
  }
  select {
  case <-ctx.Done():return nil
  case <-ticker.C:
  }
 }
}
