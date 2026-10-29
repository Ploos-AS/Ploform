// Package webhooks implements transport-independent signed event envelopes.
// Delivery, secrets storage, retries and HTTP endpoints are not implemented.
package webhooks

import (
 "crypto/hmac"
 "crypto/sha256"
 "encoding/hex"
 "errors"
 "fmt"
 "strconv"
 "strings"
 "time"
)

const (
 HeaderSignature = "X-Ploform-Signature"
 HeaderTimestamp = "X-Ploform-Timestamp"
 MaxClockSkew = 5 * time.Minute
)

// Sign authenticates the exact transmitted payload bytes and timestamp.
// Secrets must be generated and managed per organization and destination.
func Sign(secret, body []byte, timestamp time.Time) (string, error) {
 if len(secret) < 32 { return "", errors.New("webhook signing secret must be at least 32 bytes") }
 mac := hmac.New(sha256.New, secret)
 _, _ = fmt.Fprintf(mac, "%d.", timestamp.Unix())
 _, _ = mac.Write(body)
 return "v1=" + hex.EncodeToString(mac.Sum(nil)), nil
}

// Verify checks the timestamp window and signature using constant-time MAC
// comparison. Callers must also enforce event-ID deduplication.
func Verify(secret, body []byte, unixTimestamp, signature string, now time.Time) bool {
 if len(secret) < 32 || !strings.HasPrefix(signature, "v1=") { return false }
 ts, err := strconv.ParseInt(unixTimestamp, 10, 64)
 if err != nil || ts < 0 { return false }
 sent := time.Unix(ts, 0)
 if sent.Before(now.Add(-MaxClockSkew)) || sent.After(now.Add(MaxClockSkew)) { return false }
 supplied, err := hex.DecodeString(strings.TrimPrefix(signature, "v1="))
 if err != nil || len(supplied) != sha256.Size { return false }
 expected, err := Sign(secret, body, sent)
 if err != nil { return false }
 expectedBytes, err := hex.DecodeString(strings.TrimPrefix(expected, "v1="))
 return err == nil && hmac.Equal(supplied, expectedBytes)
}
