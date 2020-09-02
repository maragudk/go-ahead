package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Run("returns new server with http addresses set", func(t *testing.T) {
		s := New(Options{ExternalPort: 123, InternalPort: 234, InternalHost: "localhost"})
		require.NotNil(t, s)
		require.Equal(t, ":123", s.externalAddress)
		require.Equal(t, "localhost:234", s.internalAddress)
	})
}
