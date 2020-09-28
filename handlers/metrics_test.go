package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	t.Run("sets up a handler that writes prometheus metrics on /metrics", func(t *testing.T) {
		code, _, body := getRequest(Metrics(), "/metrics")

		require.Equal(t, http.StatusOK, code)
		require.NotEmpty(t, body)
	})
}
