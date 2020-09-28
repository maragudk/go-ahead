package handlers

import (
	"context"
	"encoding/gob"
	"net/http"

	"go-ahead/model"
)

const (
	sessionUserKey = "user"
)

var contextUserKey = struct{ name string }{sessionUserKey}

func getUserFromContext(r *http.Request) model.User {
	return r.Context().Value(contextUserKey).(model.User)
}

type signupper interface {
	Signup(ctx context.Context, name, email, password string) error
}

func Signup(repo signupper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		name := r.Form.Get("name")
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		if name == "" || email == "" || password == "" {
			http.Error(w, "name, email, or password empty", http.StatusBadRequest)
			return
		}

		if !model.ValidateEmail(email) {
			http.Error(w, "email is not an email address", http.StatusBadRequest)
			return
		}

		if err := repo.Signup(r.Context(), name, email, password); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
}

type loginner interface {
	Login(ctx context.Context, email, password string) (*model.User, error)
}

type sessionPutter interface {
	RenewToken(ctx context.Context) error
	Put(ctx context.Context, key string, value interface{})
}

func Login(repo loginner, s sessionPutter) http.HandlerFunc {
	// Register our user type with the gob encoding used by the session handler
	gob.Register(model.User{})

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email := r.Form.Get("email")
		password := r.Form.Get("password")

		if email == "" || password == "" {
			http.Error(w, "email or password empty", http.StatusBadRequest)
			return
		}

		user, err := repo.Login(r.Context(), email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if user == nil {
			http.Error(w, "email and/or password is incorrect", http.StatusForbidden)
			return
		}

		// Renew the token to avoid session fixation attacks
		if err := s.RenewToken(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Put the whole user info into the session
		s.Put(r.Context(), sessionUserKey, user)

		writeJSON(w, user)
	}
}

type sessionDestroyer interface {
	Destroy(ctx context.Context) error
}

func Logout(s sessionDestroyer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Destroy(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Session sits behind the auth middleware and just returns the current user info.
func Session() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		writeJSON(w, user)
	}
}
