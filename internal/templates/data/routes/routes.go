package routes

import (
	"{{module}}/handlers"

	"github.com/stygian91/veggies/router"
)

func InitRoutes() {
	router.Get().Group(func(g *router.Group) {
		g.HandleFunc("/", handlers.Greet)
	}).Middleware(router.LogMiddleware)
}
