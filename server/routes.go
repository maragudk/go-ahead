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

	// The REST API
	s.externalMux.Route("/api", func(r chi.Router) {
		r.Use(middleware.Recoverer, handlers.JSONHeader)
		r.Use(s.sm.LoadAndSave)

		// RPC-style handlers for authentication
		r.Post("/signup", handlers.Signup(s.Storer))
		r.Post("/login", handlers.Login(s.Storer, s.sm))
		r.Post("/logout", handlers.Logout(s.sm))

		// Handlers where the client needs to be authenticated
		r.Group(func(r chi.Router) {
			r.Use(handlers.Authorize(s.sm))

			r.Get("/session", handlers.Session())
		})
	})
}

func (s *Server) setupInternalRoutes() {
	s.internalMux.Get("/health", handlers.Health(s.Storer))
	s.internalMux.Get("/version", handlers.Version(s.Name, s.Version))
	s.internalMux.Get("/metrics", handlers.Metrics())
}
