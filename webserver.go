package hangman

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func Run() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "play", "", render.HTMLOptions{Layout: "_layout"})
	})

	m.Get("/create", func(r render.Render) {
		r.HTML(200, "create", "", render.HTMLOptions{Layout: "_layout"})
	})

	m.Run()
}
