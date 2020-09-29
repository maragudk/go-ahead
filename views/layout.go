// Package views has functions to create HTTP handlers that create views in the form of HTML pages.
package views

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/attr"
	"github.com/maragudk/gomponents/el"
)

// PageProps are properties for every Page.
type PageProps struct {
	Title string
	Path  string
	Body  g.Node
}

// Page returns a g.Node that renders an HTML document with the given props.
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
				Container(
					Header(props.Title),
					props.Body,
				),
			),
		),
	)
}

// Container restricts the width of the children, centers, and adds some padding.
func Container(children ...g.Node) g.Node {
	newChildren := []g.Node{attr.Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8")}
	newChildren = append(newChildren, children...)
	return el.Div(newChildren...)
}

// Navbar shows a navigation bar.
func Navbar(path string) g.Node {
	return el.Div(
		g.El("nav", attr.Class("bg-gray-800"),
			Container(
				el.Div(attr.Class("flex items-center justify-between h-16"),
					el.Div(attr.Class("flex items-center"),
						el.Div(attr.Class("flex-shrink-0"),
							g.El("img", attr.Class("h-8 w-8"), g.Attr("src", "/static/images/logo.svg"), g.Attr("alt", "logo")),
						),
						el.Div(attr.Class("block"),
							el.Div(attr.Class("ml-8 flex items-baseline space-x-4"),
								NavbarLink("/", "Home", path == "/"),
							),
						),
					),
					el.Div(attr.Class("flex items-center"),
						el.Div(attr.Class("block"),
							el.Div(attr.Class("flex items-baseline"),
								NavbarLink("/login", "Log in", path == "/login"),
							),
						),
					),
				),
			),
		),
	)
}

// NavbarLink is a link in the Navbar.
func NavbarLink(href, text string, active bool) g.Node {
	return g.El("a",
		g.Attr("href", href),
		attr.Classes{
			"px-3 py-2 rounded-md text-sm font-medium focus:outline-none focus:text-white focus:bg-gray-700": true,
			"text-white bg-gray-900":                           active,
			"text-gray-300 hover:text-white hover:bg-gray-700": !active,
		},
		g.Text(text),
	)
}

// Header returns a g.Node that renders a headline.
func Header(title string) g.Node {
	return g.El("header", attr.Class("bg-white"),
		Container(
			g.El("h1", attr.Class("text-3xl font-bold leading-tight text-gray-900 my-6"), g.Text(title)),
		),
	)
}
