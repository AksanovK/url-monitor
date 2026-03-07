package api

import (
	"github.com/AksanovK/url-monitor/internal/repository"
	"github.com/AksanovK/url-monitor/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AksanovK/url-monitor/internal/api/handler"
)

func NewRouter(pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repo := repository.NewMonitorRepository(pool)
	svc := service.NewMonitorService(repo)
	h := handler.NewMonitorHandler(svc)

	resultRepo := repository.NewCheckResultRepository(pool)
	resultSvc := service.NewCheckResultService(resultRepo)
	resultHandler := handler.NewCheckResultHandler(resultSvc)

	r.Route("/monitors", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Get("/{id}/results", resultHandler.List)
	})

	return r
}
