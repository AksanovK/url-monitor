package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AksanovK/url-monitor/internal/db"
	"github.com/AksanovK/url-monitor/internal/domain"
)

type CheckResultRepository struct {
	queries *db.Queries
}

func NewCheckResultRepository(pool *pgxpool.Pool) *CheckResultRepository {
	return &CheckResultRepository{
		queries: db.New(pool),
	}
}

func (r *CheckResultRepository) Save(ctx context.Context, cr *domain.CheckResult) error {
	var errText pgtype.Text
	if cr.Error != "" {
		errText = pgtype.Text{String: cr.Error, Valid: true}
	}

	return r.queries.CreateCheckResult(ctx, db.CreateCheckResultParams{
		ID:         cr.ID,
		MonitorID:  cr.MonitorID,
		StatusCode: int32(cr.StatusCode),
		LatencyMs:  int32(cr.LatencyMs),
		Error:      errText,
		CheckedAt:  cr.CheckedAt,
	})
}

func (r *CheckResultRepository) FindByMonitorWithCursor(
	ctx context.Context,
	monitorID string,
	cursor time.Time,
	limit int,
) ([]*domain.CheckResult, error) {
	rows, err := r.queries.ListCheckResultsByMonitor(ctx, db.ListCheckResultsByMonitorParams{
		MonitorID: monitorID,
		CheckedAt: cursor,
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.CheckResult, 0, len(rows))
	for _, row := range rows {
		result = append(result, toCheckResultDomain(row))
	}
	return result, nil
}

func (r *CheckResultRepository) FindLatest(ctx context.Context, monitorID string, limit int) ([]*domain.CheckResult, error) {
	rows, err := r.queries.ListLatestCheckResults(ctx, db.ListLatestCheckResultsParams{
		MonitorID: monitorID,
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*domain.CheckResult, 0, len(rows))
	for _, row := range rows {
		result = append(result, toCheckResultDomain(row))
	}
	return result, nil
}

func toCheckResultDomain(row db.CheckResult) *domain.CheckResult {
	errStr := ""
	if row.Error.Valid {
		errStr = row.Error.String
	}

	return &domain.CheckResult{
		ID:         row.ID,
		MonitorID:  row.MonitorID,
		StatusCode: int(row.StatusCode),
		LatencyMs:  int(row.LatencyMs),
		Error:      errStr,
		CheckedAt:  row.CheckedAt,
	}
}
