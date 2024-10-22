package facades

import (
	"github.com/stygian91/veggies/config"
	b "github.com/stygian91/veggies/config/base"
)

type conf struct{}

func Config() conf {
	return conf{}
}

func (this conf) App() b.App {
	appConfig, _ := config.GetGroup[b.App]("app")
	return appConfig
}
