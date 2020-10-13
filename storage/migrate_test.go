// +build integration

package storage_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/go-ahead/storagetest"
)

func TestStorer_MigrateUp(t *testing.T) {
	t.Run("migrates to newest version", func(t *testing.T) {
		s, cleanup := storagetest.CreateRootStorer()
		defer cleanup()

		// Migrate down first, because the setup code migrates up
		err := s.MigrateDown()
		require.NoError(t, err)

		err = s.MigrateUp()
		require.NoError(t, err)
	})
}

func TestStorer_MigrateDown(t *testing.T) {
	t.Run("migrates down", func(t *testing.T) {
		s, _ := storagetest.CreateRootStorer()

		err := s.MigrateDown()
		require.NoError(t, err)
	})
}
