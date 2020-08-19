package ahead

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer_versionHandler(t *testing.T) {
	t.Run("prints the version number", func(t *testing.T) {
		s := NewServer(NewServerOptions{Version: "123abc"})
		code, body := makeGETRequest(s.versionHandler(), "/version")

		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "123abc", body)
	})
}
