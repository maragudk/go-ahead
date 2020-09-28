package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"go-ahead/model"
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

type sessionGetterMock struct {
	exists bool
	user   model.User
}

func (m *sessionGetterMock) Exists(ctx context.Context, key string) bool {
	return m.exists
}

func (m *sessionGetterMock) Get(ctx context.Context, key string) interface{} {
	return m.user
}

func TestAuthorize(t *testing.T) {
	t.Run("returns 401 on no session key", func(t *testing.T) {
		repo := &sessionGetterMock{}
		h := Authorize(repo)(noopHandler())
		status, _, _ := getRequest(h, "/")
		require.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("calls next handler", func(t *testing.T) {
		repo := &sessionGetterMock{exists: true}
		called := false
		h := Authorize(repo)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		}))
		status, _, _ := getRequest(h, "/")
		require.Equal(t, http.StatusOK, status)
		require.True(t, called)
	})

	t.Run("saves user in context", func(t *testing.T) {
		repo := &sessionGetterMock{exists: true, user: model.User{
			ID: "123",
		}}
		var user model.User
		h := Authorize(repo)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user = getUserFromContext(r)
		}))
		status, _, _ := getRequest(h, "/")
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, repo.user.ID, user.ID)
	})
}
