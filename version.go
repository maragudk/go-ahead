package ahead

import (
	"io"
	"net/http"
)

func (s *Server) setupVersionHandler() {
	s.internalMux.Handle("/version", s.versionHandler())
}

func (s *Server) versionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, s.Version)
	}
}
