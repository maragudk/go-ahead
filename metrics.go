package ahead

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) setupMetricsHandler() {
	s.internalMux.Handle("/metrics", s.metricsHandler())
}

func (s *Server) metricsHandler() http.Handler {
	return promhttp.Handler()
}
