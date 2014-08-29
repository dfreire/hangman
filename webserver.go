package hangman

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type A struct {
	B B
}

type B struct {
	C int
}

func Run() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", A{B: B{C: 123}})
	})
	m.Run()
}
