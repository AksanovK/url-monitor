package api

import (
	"github.com/AksanovK/url-monitor/internal/repository"
	"github.com/AksanovK/url-monitor/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AksanovK/url-monitor/internal/api/handler"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Встроенные middleware — аналог Spring фильтров
	r.Use(middleware.Logger)    // логирует каждый запрос
	r.Use(middleware.Recoverer) // ловит panic, возвращает 500 (аналог @ExceptionHandler)

	repo := repository.NewMonitorRepository()
	svc := service.NewMonitorService(repo)
	h := handler.NewMonitorHandler(svc)

	r.Route("/monitors", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
	})

	return r
}
