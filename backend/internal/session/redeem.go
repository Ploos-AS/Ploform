package session

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

// RedeemChallenge atomically consumes a valid one-time login challenge and
// creates a session. If session creation fails, challenge consumption rolls back.
func (r Repository) RedeemChallenge(ctx context.Context, challengeToken string, lifetime time.Duration) (string,time.Time,error) {
 if r.DB==nil{return "",time.Time{},errors.New("authentication database unavailable")}
 if lifetime<=0 || lifetime>30*24*time.Hour{return "",time.Time{},errors.New("invalid session lifetime")}
 challengeDigest,err:=Digest(challengeToken)
 if err!=nil{return "",time.Time{},ErrNotFound}
 sessionToken,sessionDigest,err:=NewToken()
 if err!=nil{return "",time.Time{},err}
 tx,err:=r.DB.BeginTx(ctx,nil)
 if err!=nil{return "",time.Time{},fmt.Errorf("begin login transaction: %w",err)}
 defer tx.Rollback()
 var userID string
 err=tx.QueryRowContext(ctx,`UPDATE login_challenges
 SET consumed_at=clock_timestamp()
 WHERE token_digest=$1 AND consumed_at IS NULL
 AND expires_at>clock_timestamp() AND attempts<5
 RETURNING user_id`,challengeDigest[:]).Scan(&userID)
 if errors.Is(err,sql.ErrNoRows){return "",time.Time{},ErrNotFound}
 if err!=nil{return "",time.Time{},fmt.Errorf("consume login challenge: %w",err)}
 expiresAt:=time.Now().Add(lifetime)
 _,err=tx.ExecContext(ctx,`INSERT INTO sessions (user_id,token_digest,expires_at)
 VALUES ($1,$2,$3)`,userID,sessionDigest[:],expiresAt)
 if err!=nil{return "",time.Time{},fmt.Errorf("create login session: %w",err)}
 if err:=tx.Commit();err!=nil{return "",time.Time{},fmt.Errorf("commit login: %w",err)}
 return sessionToken,expiresAt,nil
}
