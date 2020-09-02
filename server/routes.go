package server

import (
	"go-ahead/handlers"
)

func (s *Server) setupExternalRoutes() {
	s.externalMux.Use(handlers.NoClickjacking, handlers.StrictContentSecurityPolicy)

	s.externalMux.Get("/", handlers.RootHandler(s.Storer))
}

func (s *Server) setupInternalRoutes() {
	s.internalMux.Get("/health", handlers.HealthHandler(s.Storer))
	s.internalMux.Get("/version", handlers.VersionHandler(s.Name, s.Version))
	s.internalMux.Get("/metrics", handlers.MetricsHandler())
}
