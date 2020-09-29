package views

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/el"
)

// Home creates a handler that returns the home view.
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page(PageProps{
			Title: "Home",
			Path:  r.URL.Path,
			Body:  homeBody(),
		})
		_ = g.Write(w, page)
	}
}

func homeBody() g.NodeFunc {
	return el.Div()
}
