package metrics

import (
	"net"
	"testing"
	"time"
)

func TestSession_RecordAndSnapshot(t *testing.T) {
	target := net.ParseIP("8.8.8.8")
	s := NewSession(target, 50)

	ip1 := net.ParseIP("192.168.1.1")
	ip2 := net.ParseIP("10.0.0.1")

	// Record hops
	s.Record(1, ip1, 5*time.Millisecond, true)
	s.Record(2, ip2, 15*time.Millisecond, true)
	s.Record(2, ip2, 0, false) // one loss at hop 2

	snaps := s.Snapshot()
	if len(snaps) != 2 {
		t.Fatalf("expected 2 hop snapshots, got %d", len(snaps))
	}

	if snaps[0].Sent != 1 || snaps[0].Recv != 1 {
		t.Errorf("hop1: sent=%d recv=%d, want 1/1", snaps[0].Sent, snaps[0].Recv)
	}

	if snaps[1].Sent != 2 || snaps[1].Recv != 1 {
		t.Errorf("hop2: sent=%d recv=%d, want 2/1", snaps[1].Sent, snaps[1].Recv)
	}
}

func TestSession_SetHostname(t *testing.T) {
	target := net.ParseIP("1.1.1.1")
	s := NewSession(target, 10)
	ip := net.ParseIP("192.168.0.1")

	s.Record(1, ip, 1*time.Millisecond, true)
	s.SetHostname(1, "gateway.local")

	snaps := s.Snapshot()
	if len(snaps) == 0 {
		t.Fatal("no snapshots")
	}
	if snaps[0].Hostname != "gateway.local" {
		t.Errorf("Hostname = %q, want 'gateway.local'", snaps[0].Hostname)
	}
}

func TestSession_RouteChanges(t *testing.T) {
	target := net.ParseIP("8.8.8.8")
	s := NewSession(target, 10)

	s.RecordRouteChange()
	s.RecordRouteChange()

	sum := s.Summary()
	if sum.RouteChanges != 2 {
		t.Errorf("RouteChanges = %d, want 2", sum.RouteChanges)
	}
}

func TestSession_Summary(t *testing.T) {
	target := net.ParseIP("8.8.8.8")
	s := NewSession(target, 10)
	ip := net.ParseIP("10.0.0.1")

	s.Record(1, ip, 5*time.Millisecond, true)
	s.Record(1, ip, 0, false)

	sum := s.Summary()
	if sum.TotalSent != 2 {
		t.Errorf("TotalSent = %d, want 2", sum.TotalSent)
	}
	if sum.TotalRecv != 1 {
		t.Errorf("TotalRecv = %d, want 1", sum.TotalRecv)
	}
	if !sum.Target.Equal(target) {
		t.Errorf("Target = %v, want %v", sum.Target, target)
	}
	if sum.Duration < 0 {
		t.Error("Duration should be >= 0")
	}
}
