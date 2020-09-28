package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoClickjacking(t *testing.T) {
	t.Run("sets x-frame-options and x-xss-protection headers", func(t *testing.T) {
		h := NoClickjacking(noopHandler())
		code, header, _ := getRequest(h, "/")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "deny", header.Get("X-Frame-Options"))
		require.Equal(t, "1; mode=block", header.Get("X-XSS-Protection"))
	})
}

func TestStrictContentSecurityPolicy(t *testing.T) {
	t.Run("sets the csp headers", func(t *testing.T) {
		h := StrictContentSecurityPolicy(noopHandler())
		code, header, _ := getRequest(h, "/")
		require.Equal(t, http.StatusOK, code)
		require.Equal(t, "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'", header.Get("Content-Security-Policy"))
	})
}
