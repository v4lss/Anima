package workers

import (
	"testing"
	"time"

	"github.com/v4lss/animas/internal/domain/monitor"
)

func TestQueueFor(t *testing.T) {
	tests := []struct {
		t       monitor.MonitorType
		want    string
	}{
		{monitor.HTTP, "jobs:http"},
		{monitor.HTTPS, "jobs:http"},
		{monitor.TCP, "jobs:tcp"},
	}

	for _, tt := range tests {
		t.Run(string(tt.t), func(t *testing.T) {
			got := queueFor(tt.t)
			if got != tt.want {
				t.Errorf("queueFor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDue(t *testing.T) {
	s := &Scheduler{}

	m := &monitor.Monitor{
		ID:       "test1",
		Interval: 60,
	}

	// Never checked - should be due
	if !s.isDue(m, time.Now()) {
		t.Error("monitor never checked should be due")
	}

	// Just checked - should not be due
	s.lastCheck.Store(m.ID, time.Now())
	if s.isDue(m, time.Now()) {
		t.Error("monitor just checked should not be due")
	}

	// Checked 61 seconds ago with 60s interval - should be due
	oldTime := time.Now().Add(-61 * time.Second)
	s.lastCheck.Store(m.ID, oldTime)
	if !s.isDue(m, time.Now()) {
		t.Error("monitor checked past interval should be due")
	}
}
