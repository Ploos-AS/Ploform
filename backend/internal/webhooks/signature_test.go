package webhooks

import (
 "bytes"
 "strconv"
 "testing"
 "time"
)

func TestSignatureVerification(t *testing.T) {
 secret := bytes.Repeat([]byte{0x42}, 32)
 body := []byte(`{"event":"submission.created","submission_id":"example"}`)
 now := time.Unix(1800000000, 0)
 sig, err := Sign(secret, body, now)
 if err != nil { t.Fatal(err) }
 ts := strconv.FormatInt(now.Unix(), 10)
 if !Verify(secret, body, ts, sig, now) { t.Fatal("valid signature rejected") }
 if Verify(secret, []byte("modified"), ts, sig, now) { t.Fatal("tampered body accepted") }
 if Verify(bytes.Repeat([]byte{0x43}, 32), body, ts, sig, now) { t.Fatal("wrong secret accepted") }
 if Verify(secret, body, ts, "v1=invalid", now) { t.Fatal("invalid signature accepted") }
 if Verify(secret, body, ts, sig, now.Add(6*time.Minute)) { t.Fatal("expired signature accepted") }
 if Verify(secret, body, ts, sig, now.Add(-6*time.Minute)) { t.Fatal("future signature accepted") }
 if Verify(secret, body, "not-a-timestamp", sig, now) { t.Fatal("invalid timestamp accepted") }
}

func TestRejectWeakSecret(t *testing.T) {
 if _, err := Sign([]byte("weak"), []byte("{}"), time.Now()); err == nil { t.Fatal("weak secret accepted") }
}
