package config

import "fmt"

type LoadFn func(env map[string]string) (any, error)

var (
	repository map[string]any
	loaders    map[string]LoadFn
)

func init() {
	repository = map[string]any{}
	loaders = map[string]LoadFn{}
}

func Register(name string, loader LoadFn) error {
	if _, exists := loaders[name]; exists {
		return fmt.Errorf("Config register error - group with name %s already exists", name)
	}

	loaders[name] = loader

	return nil
}

func Boot(env map[string]string) error {
	for name, loader := range loaders {
		configValue, err := loader(env)
		if err != nil {
			return fmt.Errorf("Error loading config for %s: %w", name, err)
		}

		repository[name] = configValue
	}

	return nil
}

func GetGroup[G any](name string) (G, bool) {
	if group, ok := repository[name]; ok {
		g, ok := group.(G)
		return g, ok
	}

	return *new(G), false
}
