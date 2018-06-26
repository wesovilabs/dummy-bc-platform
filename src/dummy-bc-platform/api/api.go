package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
)

// CreateRouter define router and handlers to be exposed
func CreateRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Route("/", func(r chi.Router) {
		r.Get("/", getBC)
		r.Route("/blocks", func(r chi.Router) {
			r.Post("/", createBC)
		})
	})
	return r
}
