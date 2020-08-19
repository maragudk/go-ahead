package ahead

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_healthHandler(t *testing.T) {
	t.Run("returns OK on no errors", func(t *testing.T) {
		s := NewServer(NewServerOptions{})
		code, body := makeGETRequest(s.healthHandler(), "/health")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "OK", body)
	})
}
