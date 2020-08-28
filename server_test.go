package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Run("returns new server with http ports set", func(t *testing.T) {
		s := New(Options{ExternalPort: 123, InternalPort: 234})
		require.NotNil(t, s)
		require.Equal(t, 123, s.ExternalPort)
		require.Equal(t, 234, s.InternalPort)
	})
}
