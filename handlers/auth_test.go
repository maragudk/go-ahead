package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/maragudk/go-ahead/model"
)

var now = time.Date(2020, 9, 7, 12, 0, 0, 0, time.UTC)

type signupperMock struct {
	err      error
	name     string
	email    string
	password string
}

func (m *signupperMock) Signup(ctx context.Context, name, email, password string) error {
	m.name = name
	m.email = email
	m.password = password
	return m.err
}

func TestSignup(t *testing.T) {
	t.Run("errors on empty name, email, password", func(t *testing.T) {
		repo := &signupperMock{}
		h := Signup(repo)

		status, _, _ := postRequest(h, "/signup", "name=&email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusBadRequest, status)

		status, _, _ = postRequest(h, "/signup", "name=Me&email=&password=1234567890")
		require.Equal(t, http.StatusBadRequest, status)

		status, _, _ = postRequest(h, "/signup", "name=Me&email=me%40example.com&password=")
		require.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("errors on bad email", func(t *testing.T) {
		repo := &signupperMock{}
		h := Signup(repo)

		status, _, _ := postRequest(h, "/signup", "name=Me&email=notanemail&password=1234567890")
		require.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("errors on storage error", func(t *testing.T) {
		repo := &signupperMock{err: fmt.Errorf("blerp")}
		h := Signup(repo)

		status, _, _ := postRequest(h, "/signup", "name=Me&email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusBadGateway, status)
	})

	t.Run("calls signup on repo", func(t *testing.T) {
		repo := &signupperMock{}
		h := Signup(repo)

		status, _, _ := postRequest(h, "/signup", "name=Me&email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusOK, status)

		require.Equal(t, "Me", repo.name)
		require.Equal(t, "me@example.com", repo.email)
		require.Equal(t, "1234567890", repo.password)
	})
}

type loginnerMock struct {
	err     error
	succeed bool
}

func (m *loginnerMock) Login(ctx context.Context, email, password string) (*model.User, error) {
	var user *model.User
	if m.succeed {
		user = &model.User{
			ID:       "123",
			Name:     "Me",
			Email:    "me@example.com",
			Password: "",
			Created:  now,
			Updated:  now,
		}
	}
	return user, m.err
}

type sessionPutterMock struct {
	err         error
	renewCalled bool
	key         string
	value       interface{}
}

func (m *sessionPutterMock) RenewToken(ctx context.Context) error {
	m.renewCalled = true
	return m.err
}

func (m *sessionPutterMock) Put(ctx context.Context, key string, value interface{}) {
	m.key = key
	m.value = value
}

func TestLogin(t *testing.T) {
	t.Run("errors on empty email or password", func(t *testing.T) {
		repo := &loginnerMock{}
		s := &sessionPutterMock{}
		h := Login(repo, s)

		status, _, _ := postRequest(h, "/login", "email=&password=1234567890")
		require.Equal(t, http.StatusBadRequest, status)

		status, _, _ = postRequest(h, "/login", "email=me%40example.com&password=")
		require.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("errors on storage error", func(t *testing.T) {
		repo := &loginnerMock{err: errors.New("boom")}
		s := &sessionPutterMock{}
		h := Login(repo, s)

		status, _, _ := postRequest(h, "/login", "email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusBadGateway, status)
	})

	t.Run("returns forbidden on credentials that don't match", func(t *testing.T) {
		repo := &loginnerMock{}
		s := &sessionPutterMock{}
		h := Login(repo, s)

		status, _, _ := postRequest(h, "/login", "email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusForbidden, status)
	})

	t.Run("errors on session error", func(t *testing.T) {
		repo := &loginnerMock{succeed: true}
		s := &sessionPutterMock{err: errors.New("boom")}
		h := Login(repo, s)

		status, _, _ := postRequest(h, "/login", "email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("renews session token and puts user into session", func(t *testing.T) {
		repo := &loginnerMock{succeed: true}
		s := &sessionPutterMock{}
		h := Login(repo, s)

		status, _, body := postRequest(h, "/login", "email=me%40example.com&password=1234567890")
		require.Equal(t, http.StatusOK, status)

		require.True(t, s.renewCalled)
		require.Equal(t, sessionUserKey, s.key)
		user, ok := s.value.(*model.User)
		require.True(t, ok)
		require.Equal(t, "Me", user.Name)

		require.Equal(t, `{"ID":"123","Name":"Me","Email":"me@example.com","Created":"2020-09-07T12:00:00Z","Updated":"2020-09-07T12:00:00Z"}`, body)
	})
}

type sessionDestroyerMock struct {
	err error
}

func (m *sessionDestroyerMock) Destroy(ctx context.Context) error {
	return m.err
}

func TestLogout(t *testing.T) {
	t.Run("errors on error from session", func(t *testing.T) {
		s := &sessionDestroyerMock{err: errors.New("boom")}
		h := Logout(s)

		status, _, _ := postRequest(h, "/logout", "")
		require.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("succeeds on no error from session destroy", func(t *testing.T) {
		s := &sessionDestroyerMock{}
		h := Logout(s)

		status, _, _ := postRequest(h, "/logout", "")
		require.Equal(t, http.StatusOK, status)
	})
}

func TestSession(t *testing.T) {
	t.Run("returns the user from context", func(t *testing.T) {
		user := model.User{
			ID:       "123",
			Name:     "Me",
			Email:    "me@example.com",
			Password: "",
			Created:  now,
			Updated:  now,
		}

		h := Session()

		status, _, body := request(h, http.MethodGet, "/session", context.WithValue(context.Background(), contextUserKey, user), "")
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, `{"ID":"123","Name":"Me","Email":"me@example.com","Created":"2020-09-07T12:00:00Z","Updated":"2020-09-07T12:00:00Z"}`, body)
	})
}
