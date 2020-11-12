package modules

import "github.com/tiramiseb/quickonf/output"

// Module is a module. Returns true if succeeds, false if there has been an error
type Module func(in interface{}, out output.Output, store map[string]interface{}) error

var registry = map[string]Module{}

func Register(name string, module Module) {
	registry[name] = module
}

func Get(name string) Module {
	return registry[name]
}
