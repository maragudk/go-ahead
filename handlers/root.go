package handlers

import (
	"net/http"

	g "github.com/maragudk/gomponents"

	"go-ahead/views"
)

// RootHandler handles /.
func RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := views.Page(views.PageProps{
			Title: "Home",
			Path:  r.URL.Path,
		})
		_ = g.Write(w, page)
	}
}
