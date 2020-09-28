package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"go-ahead/handlers"
	"go-ahead/views"
)

func (s *Server) setupExternalRoutes() {
	// The views that can be requested by the browser
	s.externalMux.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer, handlers.NoClickjacking, handlers.StrictContentSecurityPolicy)
		r.Use(s.sm.LoadAndSave)

		// Serve static files from the "public" directory
		staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("public")))
		r.Get("/static/*", staticHandler.ServeHTTP)

		r.Get("/", views.Home())

		r.Route("/login", func(r chi.Router) {
			r.Get("/", views.Login())
		})
	})
}

func (s *Server) setupInternalRoutes() {
	s.internalMux.Get("/health", handlers.Health(s.Storer))
	s.internalMux.Get("/version", handlers.Version(s.Name, s.Version))
	s.internalMux.Get("/metrics", handlers.Metrics())
}
