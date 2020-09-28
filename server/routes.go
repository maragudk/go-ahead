package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"go-ahead/handlers"
)

func (s *Server) setupExternalRoutes() {
	s.externalMux.Use(middleware.Recoverer, handlers.NoClickjacking, handlers.StrictContentSecurityPolicy)
	s.externalMux.Use(s.sm.LoadAndSave)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("public")))
	s.externalMux.Get("/static/*", staticHandler.ServeHTTP)

	s.externalMux.Get("/", handlers.HomeHandler())

	s.externalMux.Route("/login", func(r chi.Router) {
		r.Get("/", handlers.LoginHandler())
	})
}

func (s *Server) setupInternalRoutes() {
	s.internalMux.Get("/health", handlers.HealthHandler(s.Storer))
	s.internalMux.Get("/version", handlers.VersionHandler(s.Name, s.Version))
	s.internalMux.Get("/metrics", handlers.MetricsHandler())
}
