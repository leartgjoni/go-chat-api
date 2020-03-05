package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	// Attach router middleware.
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create API routes.
	r.Route("/", func(r chi.Router) {
		r.Get("/health", s.handlePing)

		r.HandleFunc("/ws", s.websocketHandler.Handle)
	})

	return r
}
