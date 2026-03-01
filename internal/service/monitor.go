package service

import (
	"context"
	"errors"

	"github.com/AksanovK/url-monitor/internal/domain"
	"github.com/AksanovK/url-monitor/internal/repository"
)

type MonitorService struct {
	repo *repository.MonitorRepository
}

func NewMonitorService(repo *repository.MonitorRepository) *MonitorService {
	return &MonitorService{repo: repo}
}

func (s *MonitorService) Create(ctx context.Context, url string, intervalSec int, expectedStatus int) (*domain.Monitor, error) {
	m := domain.NewMonitor(url, intervalSec, expectedStatus)

	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *MonitorService) List(ctx context.Context) ([]*domain.Monitor, error) {
	return s.repo.FindAll(ctx)
}

func (s *MonitorService) GetByID(ctx context.Context, id string) (*domain.Monitor, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("monitor not found")
	}
	return m, nil
}
