package handlers

import (
	"net/http"

	g "github.com/maragudk/gomponents"

	"go-ahead/views"
)

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := views.Page(views.PageProps{
			Title: "Login",
			Path:  r.URL.Path,
			Body:  views.Login(),
		})
		_ = g.Write(w, page)
	}
}
