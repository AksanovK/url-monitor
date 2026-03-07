package domain

import (
	"time"

	"github.com/google/uuid"
)

type CheckResult struct {
	ID         string
	MonitorID  string
	StatusCode int
	LatencyMs  int
	Error      string
	CheckedAt  time.Time
}

func NewCheckResult(monitorID string, statusCode int, latencyMs int, checkErr string) *CheckResult {
	return &CheckResult{
		ID:         uuid.New().String(),
		MonitorID:  monitorID,
		StatusCode: statusCode,
		LatencyMs:  latencyMs,
		Error:      checkErr,
		CheckedAt:  time.Now(),
	}
}
