package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Run("prints the version number", func(t *testing.T) {
		code, _, body := getRequest(Version("appy", "123abc"), "/version")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "appy 123abc", body)
	})
}
