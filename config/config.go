package config

import "fmt"

type LoadFn func(env map[string]string) (any, error)

var configRepository map[string]any

func Boot(env map[string]string, loaders map[string]LoadFn) error {
	repo := map[string]any{}

	for name, loader := range loaders {
		configValue, err := loader(env)
		if err != nil {
			return fmt.Errorf("Error loading config for %s: %w", name, err)
		}

		repo[name] = configValue
	}

	configRepository = repo

	return nil
}

func GetGroup[G any](name string) (G, bool) {
	group, ok := configRepository[name]
	if !ok {
		return *new(G), ok
	}

	g, ok := group.(G)

	return g, ok
}
