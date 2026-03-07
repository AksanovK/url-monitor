package worker

import (
	"context"
	"log"
	"sync"

	"github.com/AksanovK/url-monitor/internal/domain"
)

type Job struct {
	Monitor *domain.Monitor
}

type Result struct {
	CheckResult *domain.CheckResult
}

type Pool struct {
	numWorkers int
	checker    *Checker
	jobs       chan Job
	results    chan Result
}

func NewPool(numWorkers int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		checker:    NewChecker(),
		jobs:       make(chan Job, numWorkers*2),
		results:    make(chan Result, numWorkers*2),
	}
}

func (p *Pool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < p.numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Printf("worker %d started", id)

			for job := range p.jobs {
				select {
				case <-ctx.Done():
					log.Printf("worker %d stopped by context", id)
					return
				default:
					result := p.checker.Check(ctx, job.Monitor)
					p.results <- Result{CheckResult: result}
				}
			}

			log.Printf("worker %d finished", id)
		}(i)
	}

	go func() {
		wg.Wait()
		close(p.results)
	}()
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}

func (p *Pool) Results() <-chan Result {
	return p.results
}

func (p *Pool) Close() {
	close(p.jobs)
}
