package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler returns a new default Prometheus handler.
func MetricsHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
