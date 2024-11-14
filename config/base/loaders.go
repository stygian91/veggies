package base

import "github.com/stygian91/veggies/config"

var loaders map[string]config.LoadFn = map[string]config.LoadFn{
	"app": LoadApp,
}

func Loaders() map[string]config.LoadFn {
	return loaders
}
