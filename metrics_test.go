package ahead

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_SetupMetrics(t *testing.T) {
	t.Run("sets up a handler that writes prometheus metrics on /metrics", func(t *testing.T) {
		s := NewServer(NewServerOptions{})
		code, body := makeGETRequest(s.metricsHandler(), "/metrics")

		require.Equal(t, http.StatusOK, code)
		require.NotEmpty(t, body)
	})
}
