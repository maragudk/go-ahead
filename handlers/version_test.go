package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionHandler(t *testing.T) {
	t.Run("prints the version number", func(t *testing.T) {
		code, _, body := makeGETRequest(VersionHandler("appy", "123abc"), "/version")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "appy 123abc", body)
	})
}
