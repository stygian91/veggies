package base

import "fmt"

type App struct {
	Addr string
}

func LoadApp(env map[string]string) (any, error) {
	addr, ok := env["APP_ADDR"]
	if !ok {
		return App{}, fmt.Errorf("Error loading App config - could not find APP_ADDR")
	}

	return App{Addr: addr}, nil
}
