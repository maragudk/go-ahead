package handlers

import (
	"net/http"

	g "github.com/maragudk/gomponents"

	"go-ahead/views"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := views.Page(views.PageProps{
			Title: "Home",
			Path:  r.URL.Path,
			Body:  views.Home(),
		})
		_ = g.Write(w, page)
	}
}
