package httpauth

import (
 "errors"
 "net"
 "net/netip"
 "strings"
)

// TrustedProxyPeer validates the direct TCP peer only. It deliberately never
// consumes Forwarded or X-Forwarded-* headers.
func TrustedProxyPeer(remoteAddr, cidrs string) bool {
 host, _, err := net.SplitHostPort(remoteAddr)
 if err != nil { return false }
 peer, err := netip.ParseAddr(host)
 if err != nil { return false }
 peer = peer.Unmap()
 if peer.IsLoopback() { return true }
 if strings.TrimSpace(cidrs) == "" || ValidateTrustedProxyCIDRs(cidrs) != nil { return false }
 for _, raw := range strings.Split(cidrs, ",") {
  prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
  if err != nil { return false } // invalid config must fail closed
  if prefix.Contains(peer) { return true }
 }
 return false
}

// ValidateTrustedProxyCIDRs rejects broad or malformed proxy ranges.
// Operators must specify exact peer IPs (/32 or /128), not entire networks.
func ValidateTrustedProxyCIDRs(cidrs string) error {
 if strings.TrimSpace(cidrs) == "" { return nil }
 for _, raw := range strings.Split(cidrs, ",") {
  prefix, err := netip.ParsePrefix(strings.TrimSpace(raw))
  if err != nil { return errors.New("invalid trusted proxy CIDR") }
  addr := prefix.Addr()
  if !addr.IsValid() || addr.IsUnspecified() || addr.IsMulticast() ||
   prefix.Bits() != addr.BitLen() { return errors.New("trusted proxy must be a single IP") }
 }
 return nil
}
