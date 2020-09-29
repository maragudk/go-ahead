// +build integration

package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/go-ahead/model"
	"github.com/maragudk/go-ahead/storagetest"
)

func TestStorer_Signup(t *testing.T) {
	t.Run("creates user", func(t *testing.T) {
		s, cleanup := storagetest.CreateStorer()
		defer cleanup()

		err := s.Signup(context.Background(), "Me", "me@example.com", "1234567890")
		require.NoError(t, err)

		var user model.User
		err = s.DB.Get(&user, "select name from users where email = $1", "me@example.com")
		require.NoError(t, err)
		require.Equal(t, "Me", user.Name)
	})

	t.Run("errors on duplicate email", func(t *testing.T) {
		s, cleanup := storagetest.CreateStorer()
		defer cleanup()

		err := s.Signup(context.Background(), "Me", "me@example.com", "1234567890")
		require.NoError(t, err)

		err = s.Signup(context.Background(), "Me", "me@example.com", "1234567890")
		require.Error(t, err)
	})

	t.Run("errors on too short and too long passwords", func(t *testing.T) {
		s, cleanup := storagetest.CreateStorer()
		defer cleanup()

		err := s.Signup(context.Background(), "Me", "me@example.com", "123")
		require.Error(t, err)

		err = s.Signup(context.Background(), "Me", "me@example.com", "1234567890.1234567890.1234567890.1234567890.1234567890.1234567890")
		require.Error(t, err)
	})

	t.Run("errors on bad email", func(t *testing.T) {
		s, cleanup := storagetest.CreateStorer()
		defer cleanup()

		err := s.Signup(context.Background(), "Me", "notanemail", "1234567890")
		require.Error(t, err)
	})
}

func TestStorer_Login(t *testing.T) {
	s, cleanup := storagetest.CreateStorer()
	defer cleanup()

	err := s.Signup(context.Background(), "Me", "me@example.com", "1234567890")
	require.NoError(t, err)

	t.Run("returns no user on no email match", func(t *testing.T) {
		user, err := s.Login(context.Background(), "you@example.com", "1234567890")
		require.NoError(t, err)
		require.Nil(t, user)
	})

	t.Run("returns no user on no password match", func(t *testing.T) {
		user, err := s.Login(context.Background(), "me@example.com", "2345678901")
		require.NoError(t, err)
		require.Nil(t, user)
	})

	t.Run("returns user on email and password match, without the password", func(t *testing.T) {
		user, err := s.Login(context.Background(), "me@example.com", "1234567890")
		require.NoError(t, err)
		require.NotNil(t, user)

		require.Equal(t, "Me", user.Name)
		require.Equal(t, "me@example.com", user.Email)
		require.Equal(t, "", user.Password)
	})

	t.Run("returns no user on no email match and empty passwords", func(t *testing.T) {
		user, err := s.Login(context.Background(), "you@example.com", "")
		require.NoError(t, err)
		require.Nil(t, user)
	})
}
