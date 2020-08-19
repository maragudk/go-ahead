package ahead

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	s := NewServer(NewServerOptions{})
	s.setupHealthHandler()

	t.Run("returns OK on no errors", func(t *testing.T) {
		code, body := makeGETRequest(s.healthHandler(), "/health")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "OK", body)
	})
}
