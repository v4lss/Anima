package monitor

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		monitor *Monitor
		wantErr bool
	}{
		{
			name: "valid monitor",
			monitor: &Monitor{
				Name:     "Test API",
				Target:   "https://api.example.com",
				Type:     HTTPS,
				Interval: 60,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			monitor: &Monitor{
				Name:     "",
				Target:   "https://api.example.com",
				Type:     HTTPS,
				Interval: 60,
			},
			wantErr: true,
		},
		{
			name: "empty target",
			monitor: &Monitor{
				Name:     "Test API",
				Target:   "",
				Type:     HTTPS,
				Interval: 60,
			},
			wantErr: true,
		},
		{
			name: "interval too low",
			monitor: &Monitor{
				Name:     "Test API",
				Target:   "https://api.example.com",
				Type:     HTTPS,
				Interval: 10,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			monitor: &Monitor{
				Name:     "Test API",
				Target:   "https://api.example.com",
				Type:     MonitorType("INVALID"),
				Interval: 60,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.monitor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
