package views

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

type PageProps struct {
	Title string
	Path  string
}

func Page(props PageProps) g.Node {
	return el.Document(
		el.HTML(g.Attr("lang", "en"),
			el.Head(
				el.Meta(g.Attr("charset", "utf-8")),
				el.Meta(g.Attr("name", "viewport"), g.Attr("content", "width=device-width, initial-scale=1")),
				el.Title(props.Title),
				el.Link(g.Attr("rel", "stylesheet"), g.Attr("href", "/static/styles/app.css")),
			),
			el.Body(
				Navbar(props.Path),
				Header(props.Title),
			),
		),
	)
}

func Navbar(path string) g.Node {
	return el.Div(
		g.El("nav", attr.Class("bg-gray-800"),
			el.Div(attr.Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
				el.Div(attr.Class("flex items-center justify-between h-16"),
					el.Div(attr.Class("flex items-center"),
						el.Div(attr.Class("flex-shrink-0"),
							g.El("img", attr.Class("h-8 w-8"), g.Attr("src", "/static/images/logo.svg"), g.Attr("alt", "logo")),
						),
						el.Div(attr.Class("block"),
							el.Div(attr.Class("ml-10 flex items-baseline space-x-4"),
								NavbarLink("/", "Home", path == "/"),
								NavbarLink("#", "Page 1", path == "/page1"),
								NavbarLink("#", "Page 2", path == "/page2"),
							),
						),
					),
				),
			),
		),
	)
}

func NavbarLink(href, text string, active bool) g.Node {
	return g.El("a",
		g.Attr("href", href),
		attr.Classes{
			"px-3 py-2 rounded-md text-sm font-medium":                                                               true,
			"text-white bg-gray-900 focus:outline-none focus:text-white focus:bg-gray-700":                           active,
			"text-gray-300 hover:text-white hover:bg-gray-700 focus:outline-none focus:text-white focus:bg-gray-700": !active,
		},
		g.Text(text),
	)
}

func Header(title string) g.Node {
	return g.El("header", attr.Class("bg-white shadow"),
		el.Div(attr.Class("max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8"),
			g.El("h1", attr.Class("text-3xl font-bold leading-tight text-gray-900"), g.Text(title)),
		),
	)
}
