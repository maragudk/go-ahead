package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type pingerMock struct {
	called bool
}

func (p *pingerMock) Ping(ctx context.Context) error {
	p.called = true
	return nil
}

func TestHealth(t *testing.T) {
	t.Run("returns OK on no errors", func(t *testing.T) {
		p := &pingerMock{}
		code, _, body := getRequest(Health(p), "/health")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "OK", body)
		require.True(t, p.called)
	})
}
