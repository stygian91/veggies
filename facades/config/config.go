package config

import (
	"github.com/stygian91/veggies/config"
	b "github.com/stygian91/veggies/config/base"
)

func App() b.App {
	appConfig, _ := config.GetGroup[b.App]("app")
	return appConfig
}

func Get[G any](name string) (G, bool) {
	return config.GetGroup[G](name)
}

func Register(name string, loader config.LoadFn) error {
	return config.Register(name, loader)
}

func Env() map[string]string {
	return config.GetEnv()
}
