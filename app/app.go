package app

import (
	"github.com/stygian91/veggies/config"
	"github.com/stygian91/veggies/config/base"
	"github.com/stygian91/veggies/env"
)

func Boot() error {
	if err := env.Boot(); err != nil {
		return err
	}

	if err := registerBaseConfig(); err != nil {
		return err
	}

	if err := config.Boot(env.Get()); err != nil {
		return err
	}

	return nil
}

func Run() error {
	// TODO:
	return nil
}

func registerBaseConfig() error {
	for name, loader := range base.Loaders() {
		if err := config.Register(name, loader); err != nil {
			return err
		}
	}

	return nil
}
