package service

import (
	"context"
	"time"

	"github.com/AksanovK/url-monitor/internal/domain"
	"github.com/AksanovK/url-monitor/internal/repository"
)

type CheckResultService struct {
	repo *repository.CheckResultRepository
}

func NewCheckResultService(repo *repository.CheckResultRepository) *CheckResultService {
	return &CheckResultService{repo: repo}
}

func (s *CheckResultService) GetByMonitor(ctx context.Context, monitorID string, cursor *time.Time, limit int) ([]*domain.CheckResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if cursor != nil {
		return s.repo.FindByMonitorWithCursor(ctx, monitorID, *cursor, limit)
	}
	return s.repo.FindLatest(ctx, monitorID, limit)
}
