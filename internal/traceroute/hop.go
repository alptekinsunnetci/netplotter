// Package traceroute discovers and tracks the network path to a target.
package traceroute

import (
	"net"
	"time"
)

// HopState describes what was last seen at a TTL position.
type HopState int

const (
	HopUnknown     HopState = iota // no reply received yet
	HopIntermediate                 // received TTL Exceeded
	HopDestination                  // received Echo Reply / TCP RST / etc.
	HopNoReply                      // consistently timed out
)

// Hop represents one entry in a discovered network path.
type Hop struct {
	TTL        int
	IP         net.IP
	Hostname   string   // reverse DNS (may be empty)
	State      HopState
	DiscoveredAt time.Time
	LastSeen   time.Time
}

// Equal returns true if both hops have the same TTL and responding IP.
func (h *Hop) Equal(other *Hop) bool {
	if h == nil || other == nil {
		return h == other
	}
	return h.TTL == other.TTL && h.IP.Equal(other.IP)
}

// DisplayName returns the hostname if available, otherwise the IP string.
func (h *Hop) DisplayName() string {
	if h.Hostname != "" {
		return h.Hostname
	}
	if h.IP != nil {
		return h.IP.String()
	}
	return "???"
}
