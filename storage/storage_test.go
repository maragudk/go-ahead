// +build integration

package storage_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go-ahead/storagetest"
)

func TestMain(m *testing.M) {
	storagetest.HandleTestMain(m)
}

func TestStorer_Ping(t *testing.T) {
	t.Run("pings with no error", func(t *testing.T) {
		s, cleanup := storagetest.CreateStorer()
		defer cleanup()

		err := s.Ping()
		require.NoError(t, err)
	})
}
