package notify

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

// Job is a tenant-scoped queued email. Address and body are never supplied by
// an unauthenticated form request; enqueue only from trusted server configuration.
type Job struct {
 ID int64
 OrganizationID string
 SubmissionID string
 Recipient string
 Subject string
 Body string
 Attempts int
}

type Outbox struct { DB *sql.DB }

// EnqueueTx must be called inside the same transaction that stores the submission.
// The unique key prevents duplicate notifications for one submission.
func (o Outbox) EnqueueTx(ctx context.Context, tx *sql.Tx, j Job) error {
 if tx==nil || j.OrganizationID=="" || j.SubmissionID=="" || j.Recipient=="" {
  return errors.New("invalid notification job")
 }
 _,err:=tx.ExecContext(ctx,`INSERT INTO notification_outbox
  (organization_id,submission_id,recipient,subject,body)
  VALUES ($1,$2,$3,$4,$5)
  ON CONFLICT (organization_id,submission_id) DO NOTHING`,
  j.OrganizationID,j.SubmissionID,j.Recipient,j.Subject,j.Body)
 return err
}

// Claim atomically leases one due job. SKIP LOCKED allows multiple workers.
func (o Outbox) Claim(ctx context.Context) (Job,error) {
 var j Job
 if o.DB==nil{return j,errors.New("database required")}
 err:=o.DB.QueryRowContext(ctx,`WITH next AS (
 SELECT id FROM notification_outbox
 WHERE status IN ('pending','retry') AND next_attempt_at <= now()
    OR (status='sending' AND lease_until < now())
 ORDER BY next_attempt_at,id
 FOR UPDATE SKIP LOCKED LIMIT 1
 )
 UPDATE notification_outbox n SET status='sending', attempts=n.attempts+1,
 lease_until=now()+interval '2 minutes', updated_at=now()
 FROM next WHERE n.id=next.id
 RETURNING n.id,n.organization_id::text,n.submission_id::text,n.recipient,n.subject,n.body,n.attempts`).Scan(
 &j.ID,&j.OrganizationID,&j.SubmissionID,&j.Recipient,&j.Subject,&j.Body,&j.Attempts)
 return j,err
}

// Complete uses the attempt count as a fencing token to avoid stale workers
// overwriting a later attempt's state.
func (o Outbox) Complete(ctx context.Context,j Job,sendErr error) error {
 if o.DB==nil{return errors.New("database required")}
 if sendErr==nil {
  _,err:=o.DB.ExecContext(ctx,`UPDATE notification_outbox SET status='sent',
   lease_until=NULL, sent_at=now(), updated_at=now(), body=''
   WHERE id=$1 AND attempts=$2 AND status='sending'`,j.ID,j.Attempts)
  return err
 }
 delay:=time.Duration(1<<min(j.Attempts,8))*time.Minute
 status:="retry"
 if j.Attempts>=10 {status="failed"}
 _,err:=o.DB.ExecContext(ctx,`UPDATE notification_outbox SET status=$3,
 lease_until=NULL,next_attempt_at=now()+($4 * interval '1 second'),updated_at=now()
 WHERE id=$1 AND attempts=$2 AND status='sending'`,
 j.ID,j.Attempts,status,int(delay.Seconds()))
 return err
}

func (o Outbox) ProcessOne(ctx context.Context, send func(context.Context,Job) error) error {
 if send==nil{return errors.New("sender required")}
 j,err:=o.Claim(ctx)
 if errors.Is(err,sql.ErrNoRows){return nil}
 if err!=nil{return fmt.Errorf("claim notification: %w",err)}
 sendErr:=send(ctx,j)
 if err:=o.Complete(ctx,j,sendErr);err!=nil{return fmt.Errorf("complete notification: %w",err)}
 // Delivery failures are represented as queued retry state, not logged with PII.
 return nil
}
