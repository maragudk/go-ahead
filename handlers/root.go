package handlers

import (
	"context"
	"net/http"
)

type pinger interface {
	Ping(context.Context) error
}

// RootHandler handles everything on the external interface.
func RootHandler(p pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := p.Ping(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
}
