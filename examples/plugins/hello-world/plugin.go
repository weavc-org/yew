package main

import (
	"github.com/weavc/yew/v3/pkg"
)

// Plugin variable must be exported
// Plugin must also implement the github.com/weavc/yew/v3/pkg/plugin.Plugin interface
var Plugin HelloWorldPlugin = HelloWorldPlugin{}

type HelloWorldPlugin struct {
	handler pkg.Handler
	pkg.Plugin
}

type Config struct {
	Say string
}

// Manifest returns
func (p *HelloWorldPlugin) Manifest() pkg.Manifest {
	return pkg.Manifest{
		Namespace:   "examples.helloworld",
		Description: "Hello world event plugin",
	}
}

// Register is used to initialize & setup the plugin
func (p *HelloWorldPlugin) Setup(handler pkg.Handler) error {
	p.handler = handler
	return nil
}
