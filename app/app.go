package app

import (
	"fmt"
	"net/http"

	"github.com/stygian91/veggies/config"
	"github.com/stygian91/veggies/config/base"
	"github.com/stygian91/veggies/env"
	"github.com/stygian91/veggies/router"
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

	router.Get().Boot()

	return nil
}

func Run() error {
	appCfg, ok := config.GetGroup[base.App]("app")
	if !ok {
		return fmt.Errorf("Error while starting up app - invalid app config.")
	}

	fmt.Printf("Starting server on: %s\n", appCfg.Addr)

	// TODO: check for SSL config and use ListenAndServeTLS if it's available
	return http.ListenAndServe(appCfg.Addr, router.Get().Mux())
}

func registerBaseConfig() error {
	for name, loader := range base.Loaders() {
		if err := config.Register(name, loader); err != nil {
			return err
		}
	}

	return nil
}
