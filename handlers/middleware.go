package handlers

import (
	"context"
	"net/http"

	"github.com/maragudk/go-ahead/model"
)

// Middleware is an alias for a function that takes a handler and returns one, too.
type Middleware = func(http.Handler) http.Handler

// NoClickjacking middleware sets headers to disallow frame embedding and XSS protection for older browsers.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
func NoClickjacking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

// StrictContentSecurityPolicy sets best practice CSP headers.
// This disallows all external img, script, and style links, and disallows all objects (flash etc.).
// See https://infosec.mozilla.org/guidelines/web_security#content-security-policy
func StrictContentSecurityPolicy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'")
		next.ServeHTTP(w, r)
	})
}

// JSONHeader adds a JSON Content-Type header.
func JSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type sessionGetter interface {
	Exists(ctx context.Context, key string) bool
	Get(ctx context.Context, key string) interface{}
}

// Authorize checks that there's a user logged in, and otherwise returns HTTP 401 Unauthorized.
func Authorize(s sessionGetter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !s.Exists(r.Context(), sessionUserKey) {
				http.Error(w, "unauthorized, please login", http.StatusUnauthorized)
				return
			}

			user, ok := s.Get(r.Context(), sessionUserKey).(model.User)
			if !ok {
				http.Error(w, "could not hydrate user", http.StatusInternalServerError)
				return
			}

			// We store the user directly in the context instead of having to use the session manager
			ctx := context.WithValue(r.Context(), contextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
