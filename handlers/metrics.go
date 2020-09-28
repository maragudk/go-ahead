package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics returns a new default Prometheus handler.
func Metrics() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
