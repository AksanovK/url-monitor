package worker

import (
	"context"
	"log"
	"time"

	"github.com/AksanovK/url-monitor/internal/repository"
)

type Scheduler struct {
	monitorRepo *repository.MonitorRepository
	resultRepo  *repository.CheckResultRepository
	pool        *Pool
	interval    time.Duration
}

func NewScheduler(
	monitorRepo *repository.MonitorRepository,
	resultRepo *repository.CheckResultRepository,
	numWorkers int,
	interval time.Duration,
) *Scheduler {
	return &Scheduler{
		monitorRepo: monitorRepo,
		resultRepo:  resultRepo,
		pool:        NewPool(numWorkers),
		interval:    interval,
	}
}

func (s *Scheduler) submitChecks(ctx context.Context) {
	monitors, err := s.monitorRepo.FindAll(ctx)

	if err != nil {
		log.Printf("failed to fetch monitors: %v", err)
		return
	}

	log.Printf("submitting %d checks", len(monitors))

	for _, monitor := range monitors {
		s.pool.Submit(Job{Monitor: monitor})
	}
}

func (s *Scheduler) scheduleChecks(ctx context.Context) {
	s.submitChecks(ctx)
	ticker := time.NewTicker(s.interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.pool.Close()
			log.Println("scheduler stopped")
			return
		case <-ticker.C:
			s.submitChecks(ctx)
		}
	}

}

func (s *Scheduler) collectResults(ctx context.Context) {
	for result := range s.pool.Results() {
		cr := result.CheckResult

		if cr.Error != "" {
			log.Printf("CHECK FAIL monitor=%s latency=%dms error=%s", cr.MonitorID, cr.LatencyMs, cr.Error)
		} else {
			log.Printf("CHECK OK monitor=%s status=%d latency=%dms", cr.MonitorID, cr.StatusCode, cr.LatencyMs)
		}

		if err := s.resultRepo.Save(ctx, cr); err != nil {
			log.Printf("failed to save check result: %v", err)
		}
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	s.pool.Start(ctx)

	go s.scheduleChecks(ctx)
	go s.collectResults(ctx)

	log.Printf("scheduler started: %d workers, interval %s", s.pool.numWorkers, s.interval)
}
