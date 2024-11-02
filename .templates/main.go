package main

import (
	"veggie-app/routes"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stygian91/veggies/app"
)

func main() {
	routes.InitRoutes()

	app.Boot()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
