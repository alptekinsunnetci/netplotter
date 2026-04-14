//go:build windows

package probe

import (
	"fmt"

	"github.com/alptekinsunnetci/netplotter/internal/config"
)

// New creates a Prober.  On Windows we always prefer WindowsICMPProber which
// uses IcmpSendEcho2 — this is the only way to receive ICMP Time Exceeded on
// Windows, because the kernel blocks type-11 from reaching raw sockets.
func New(cfg *config.Config) (Prober, error) {
	switch cfg.Protocol {
	case config.ProtoTCP:
		return NewTCPProber(cfg.Port), nil

	default: // icmp, udp, or anything else → Windows ICMP API
		p, err := NewWindowsICMPProber()
		if err != nil {
			fmt.Printf("[warn] IcmpSendEcho2 unavailable (%v), falling back to TCP/%d\n", err, cfg.Port)
			return NewTCPProber(cfg.Port), nil
		}
		return p, nil
	}
}
