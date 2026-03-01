package domain

import "testing"

func TestMonitorValidate(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		intervalSec int
		status      int
		wantErr     bool
	}{
		{
			name:        "valid monitor",
			url:         "https://google.com",
			intervalSec: 30,
			status:      200,
			wantErr:     false,
		},
		{
			name:        "empty url",
			url:         "",
			intervalSec: 30,
			status:      200,
			wantErr:     true,
		},
		{
			name:        "zero interval",
			url:         "https://google.com",
			intervalSec: 0,
			status:      200,
			wantErr:     true,
		},
		{
			name:        "negative interval",
			url:         "https://google.com",
			intervalSec: -5,
			status:      200,
			wantErr:     true,
		},
		{
			name:        "invalid status low",
			url:         "https://google.com",
			intervalSec: 30,
			status:      99,
			wantErr:     true,
		},
		{
			name:        "invalid status high",
			url:         "https://google.com",
			intervalSec: 30,
			status:      600,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMonitor(tt.url, tt.intervalSec, tt.status)
			err := m.Validate()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

func TestNewMonitor(t *testing.T) {
	m := NewMonitor("https://google.com", 30, 200)

	if m.ID == "" {
		t.Error("expected ID to be generated")
	}
	if m.URL != "https://google.com" {
		t.Errorf("expected url https://google.com, got %s", m.URL)
	}
	if m.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
