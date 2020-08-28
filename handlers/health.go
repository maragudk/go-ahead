package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
)

// HealthHandler for load balancers to check if this node is alive.
// It can also print a log line to stderr.
func HealthHandler(p pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if _, ok := q["logError"]; ok {
			l := log.New(os.Stderr, "", 0)
			l.Println("Testing error log in health endpoint")
		}

		if err := p.Ping(r.Context()); err != nil {
			http.Error(w, "could not ping", http.StatusBadGateway)
			return
		}

		_, _ = io.WriteString(w, "OK")
	}
}