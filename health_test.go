package ahead

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	s := NewServer(NewServerOptions{})
	s.setupHealthHandler()

	t.Run("returns OK on no errors", func(t *testing.T) {
		request := httptest.NewRequest("GET", "/health", nil)
		recorder := httptest.NewRecorder()

		s.internalMux.ServeHTTP(recorder, request)
		result := recorder.Result()

		body, err := ioutil.ReadAll(result.Body)
		require.NoError(t, err)

		require.Equal(t, "OK", string(body))
	})
}
