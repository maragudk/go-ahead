package server

import (
	"net/http"

	"go-ahead/handlers"
)

func (s *Server) setupExternalRoutes() {
	s.externalMux.Use(handlers.NoClickjacking, handlers.StrictContentSecurityPolicy)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("public")))
	s.externalMux.Get("/static/*", staticHandler.ServeHTTP)

	s.externalMux.Get("/", handlers.HomeHandler())
	s.externalMux.Get("/login", handlers.LoginHandler())
}

func (s *Server) setupInternalRoutes() {
	s.internalMux.Get("/health", handlers.HealthHandler(s.Storer))
	s.internalMux.Get("/version", handlers.VersionHandler(s.Name, s.Version))
	s.internalMux.Get("/metrics", handlers.MetricsHandler())
}
