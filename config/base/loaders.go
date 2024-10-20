package base

import "github.com/stygian91/veggies/config"

var loaders map[string]config.LoadFn

func init() {
	ls := map[string]config.LoadFn{
		"app": LoadApp,
	}

	loaders = ls
}

func Loaders() map[string]config.LoadFn {
	return loaders
}
