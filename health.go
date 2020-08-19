package ahead

import (
	"io"
	"log"
	"net/http"
	"os"
)

func (s *Server) setupHealthHandler() {
	s.internalMux.Handle("/health", s.healthHandler())
}

// healthHandler for load balancers to check if this node is alive.
// It can also print a log line to stderr.
func (s *Server) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if _, ok := q["logError"]; ok {
			l := log.New(os.Stderr, "", 0)
			l.Println("Testing error log in health endpoint")
		}

		_, _ = io.WriteString(w, "OK")
	}
}
