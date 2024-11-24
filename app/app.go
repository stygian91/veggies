package app

import (
	"fmt"
	"net/http"

	"github.com/stygian91/veggies/config"
	"github.com/stygian91/veggies/config/base"
	"github.com/stygian91/veggies/router"
)

func Boot() error {
	if err := config.BootEnv(); err != nil {
		return err
	}

	if err := registerBaseConfig(); err != nil {
		return err
	}

	if err := config.Boot(config.GetEnv()); err != nil {
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

	if len(appCfg.SSLCert) > 0 && len(appCfg.SSLKey) > 0 {
		fmt.Printf("Starting https server on: %s\n", appCfg.Addr)
		return http.ListenAndServeTLS(appCfg.Addr, appCfg.SSLCert, appCfg.SSLKey, router.Get().Mux())
	}

	fmt.Printf("Starting http server on: %s\n", appCfg.Addr)
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
