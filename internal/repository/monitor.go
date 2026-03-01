package repository

import (
	"sync"

	"github.com/AksanovK/url-monitor/internal/domain"
)

type MonitorRepository struct {
	mu       sync.RWMutex
	monitors map[string]*domain.Monitor
}

func NewMonitorRepository() *MonitorRepository {
	return &MonitorRepository{
		monitors: make(map[string]*domain.Monitor),
	}
}

func (r *MonitorRepository) Save(m *domain.Monitor) error {
	r.mu.Lock()
	r.monitors[m.ID] = m
	r.mu.Unlock()
	return nil
}

func (r *MonitorRepository) FindAll() ([]*domain.Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Monitor, 0, len(r.monitors))
	for _, m := range r.monitors {
		result = append(result, m)
	}
	return result, nil
}

func (r *MonitorRepository) FindByID(id string) (*domain.Monitor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.monitors[id]
	if !ok {
		return nil, nil
	}
	return m, nil
}
