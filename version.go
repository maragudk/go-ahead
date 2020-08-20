package ahead

import (
	"fmt"
	"io"
	"net/http"
)

func (s *Server) setupVersionHandler() {
	s.internalMux.Handle("/version", s.versionHandler())
}

func (s *Server) versionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, fmt.Sprintf("%v %v", s.Name, s.Version))
	}
}
