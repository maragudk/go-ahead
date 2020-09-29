package views

import (
	"net/http"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

// Login creates a handler that returns the login view.
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := Page(PageProps{
			Title: "Log in",
			Path:  r.URL.Path,
			Body:  loginBody(),
		})
		_ = g.Write(w, page)
	}
}

func loginBody() g.NodeFunc {
	return el.Div(
		el.Form("/login", "post", attr.Class("max-w-sm mx-auto space-y-6"),
			el.Div(
				el.Label("email", g.Text("Email"), attr.Class("block")),
				el.Input("email", "email", attr.ID("email"), attr.Required(), attr.Placeholder("me@example.com"),
					attr.Class("form-input block w-full rounded-md shadow-sm mt-1")),
			),
			el.Div(
				el.Label("password", g.Text("Password"), attr.Class("block")),
				el.Div(attr.Class("rounded-md shadow-sm"),
					el.Input("password", "password", attr.ID("password"), attr.Required(), attr.Placeholder("******"),
						attr.Class("form-input block w-full rounded-md shadow-sm mt-1")),
				),
			),
			el.Input("submit", "submit", g.Attr("value", "Log in"),
				attr.Class("px-3 py-2 rounded-md text-lg text-white bg-gray-900 hover:bg-gray-700 cursor-pointer"),
			),
		),
	)
}
