package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID             string
	URL            string
	IntervalSec    int
	ExpectedStatus int
	CreatedAt      time.Time
}

func NewMonitor(url string, intervalSec int, expectedStatus int) *Monitor {
	return &Monitor{
		ID:             uuid.New().String(),
		URL:            url,
		IntervalSec:    intervalSec,
		ExpectedStatus: expectedStatus,
		CreatedAt:      time.Now(),
	}
}

func (m *Monitor) Validate() error {
	if m.URL == "" {
		return errors.New("url is required")
	}
	if m.IntervalSec <= 0 {
		return errors.New("interval must be positive")
	}
	if m.ExpectedStatus < 100 || m.ExpectedStatus > 599 {
		return errors.New("expected status must be between 100 and 599")
	}
	return nil
}
