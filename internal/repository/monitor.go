package repository

import (
	"context"

	"github.com/AksanovK/url-monitor/internal/db"
	"github.com/AksanovK/url-monitor/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitorRepository struct {
	queries *db.Queries
}

func NewMonitorRepository(pool *pgxpool.Pool) *MonitorRepository {
	return &MonitorRepository{
		queries: db.New(pool),
	}
}

func (r *MonitorRepository) Save(ctx context.Context, m *domain.Monitor) error {
	return r.queries.CreateMonitor(ctx, db.CreateMonitorParams{
		ID:             m.ID,
		Url:            m.URL,
		IntervalSec:    int32(m.IntervalSec),
		ExpectedStatus: int32(m.ExpectedStatus),
		CreatedAt:      m.CreatedAt,
	})
}

func (r *MonitorRepository) FindAll(ctx context.Context) ([]*domain.Monitor, error) {
	rows, err := r.queries.ListMonitors(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Monitor, 0, len(rows))
	for _, row := range rows {
		result = append(result, toDomain(row))
	}
	return result, nil
}

func (r *MonitorRepository) FindByID(ctx context.Context, id string) (*domain.Monitor, error) {
	row, err := r.queries.GetMonitorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomain(row), nil
}

func toDomain(row db.Monitor) *domain.Monitor {
	return &domain.Monitor{
		ID:             row.ID,
		URL:            row.Url,
		IntervalSec:    int(row.IntervalSec),
		ExpectedStatus: int(row.ExpectedStatus),
		CreatedAt:      row.CreatedAt,
	}
}
