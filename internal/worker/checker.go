package worker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AksanovK/url-monitor/internal/domain"
)

type Checker struct {
	client *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Checker) Check(ctx context.Context, m *domain.Monitor) *domain.CheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return domain.NewCheckResult(m.ID, 0, 0, fmt.Sprintf("create request error: %v", err))
	}

	resp, err := c.client.Do(req)
	latencyMs := int(time.Since(start).Milliseconds())

	if err != nil {
		return domain.NewCheckResult(m.ID, 0, latencyMs, fmt.Sprintf("request error: %v", err))
	}
	defer resp.Body.Close()

	return domain.NewCheckResult(m.ID, resp.StatusCode, latencyMs, "")
}
