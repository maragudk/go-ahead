package errors2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	t.Run("returns wrapped error", func(t *testing.T) {
		err := Wrap(errors.New("oops"), "wrap this")
		require.Error(t, err)
		require.Equal(t, "wrap this: oops", err.Error())
		unwrappedErr := errors.Unwrap(err)
		require.Error(t, unwrappedErr)
		require.Equal(t, "oops", unwrappedErr.Error())
	})

	t.Run("can use format parameters", func(t *testing.T) {
		err := Wrap(errors.New("whoops"), "error in %v", 123)
		require.Error(t, err)
		require.Equal(t, "error in 123: whoops", err.Error())
	})
}
