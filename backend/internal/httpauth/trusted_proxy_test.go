package httpauth

import "testing"

func TestTrustedProxyPeer(t *testing.T) {
 for _, tc := range []struct { peer, ips string; want bool }{
  {"127.0.0.1:8000", "", true},
  {"[::1]:8000", "", true},
  {"172.18.0.2:54321", "", false},
  {"172.18.0.2:54321", "172.18.0.2/32", true},
  {"172.18.0.3:54321", "172.18.0.2/32", false},
  {"172.18.0.2:54321", "invalid", false},
  {"172.18.0.2", "172.18.0.2/32", false},
 } {
  if got := TrustedProxyPeer(tc.peer, tc.ips); got != tc.want {
   t.Errorf("TrustedProxyPeer(%q,%q)=%v want %v", tc.peer, tc.ips, got, tc.want)
  }
 }
}
func TestValidateTrustedProxyCIDRs(t *testing.T) {
 for _, good := range []string{"", "172.18.0.2/32", "2001:db8::1/128", "172.18.0.2/32, 172.18.0.3/32"} {
  if err := ValidateTrustedProxyCIDRs(good); err != nil { t.Errorf("%q: %v", good, err) }
 }
 for _, bad := range []string{"0.0.0.0/0", "172.18.0.0/24", "invalid", "172.18.0.2", "172.18.0.2/32,"} {
  if err := ValidateTrustedProxyCIDRs(bad); err == nil { t.Errorf("%q accepted", bad) }
 }
}
