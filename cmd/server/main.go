package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AksanovK/url-monitor/internal/repository"
	"github.com/AksanovK/url-monitor/internal/worker"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AksanovK/url-monitor/internal/api"
	"github.com/AksanovK/url-monitor/internal/config"
)

func main() {
	cfg := config.Load()

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create migrate: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	monitorRepo := repository.NewMonitorRepository(pool)
	resultRepo := repository.NewCheckResultRepository(pool)
	scheduler := worker.NewScheduler(monitorRepo, resultRepo, 5, 30*time.Second)
	scheduler.Start(ctx)

	router := api.NewRouter(pool)
	server := &http.Server{Addr: cfg.Addr(), Handler: router}

	go func() {
		fmt.Printf("Server starting on %s\n", cfg.Addr())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}
