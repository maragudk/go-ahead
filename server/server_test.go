package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Run("returns new server with http addresses set", func(t *testing.T) {
		s := New(Options{Port: 123})
		require.NotNil(t, s)
		require.Equal(t, ":123", s.address)
	})
}
