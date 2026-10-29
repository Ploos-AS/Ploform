package session

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

// ChallengeIssuer uses a dedicated credential authorized to issue login
// challenges. Do not connect it using the public application role.
type ChallengeIssuer struct { DB *sql.DB }

// IssueChallenge generates a single-use bearer token for an existing user.
// Delivery, identity lookup and rate limiting must be enforced by its caller.
// No HTTP route invokes this method until those safeguards are implemented.
func (i ChallengeIssuer) IssueChallenge(ctx context.Context, userID string) (string,error) {
 if i.DB==nil {return "",errors.New("challenge issuer database unavailable")}
 if userID=="" {return "",errors.New("user required")}
 token,digest,err:=NewToken()
 if err!=nil{return "",err}
 _,err=i.DB.ExecContext(ctx,`INSERT INTO login_challenges
 (user_id,token_digest,expires_at)
 VALUES ($1,$2,clock_timestamp()+interval '10 minutes')`,userID,digest[:])
 if err!=nil{return "",fmt.Errorf("issue login challenge: %w",err)}
 return token,nil
}

// ChallengeLifetime documents the intended lifetime of issued credentials.
const ChallengeLifetime = 10*time.Minute
