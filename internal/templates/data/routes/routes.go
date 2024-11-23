package routes

import (
	"{{module}}/handlers"

	"github.com/stygian91/veggies/router"
	m "github.com/stygian91/veggies/router/middleware"
)

func InitRoutes() {
	router.Get().Group(func(g *router.Group) {
		g.HandleFunc("/", handlers.Greet)
	}).Middleware(m.LogMiddleware)
}
