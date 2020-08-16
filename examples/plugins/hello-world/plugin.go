package main

import (
	"log"

	"github.com/weavc/yew/pkg"
)

// Plugin variable must be exported
// Plugin must also implement the github.com/weavc/yew/pkg/plugin.Plugin interface
var Plugin HelloWorldPlugin = HelloWorldPlugin{}

type HelloWorldPlugin struct {
	handler pkg.Handler
	pkg.Plugin
}

// Manifest returns
func (p *HelloWorldPlugin) Manifest() pkg.Manifest {
	return pkg.Manifest{
		Namespace:   "examples.helloworld",
		Description: "Hello world event plugin",
		Events:      map[string]func(event string, v interface{}){pkg.LoadedEvent: helloWorldEvent},
	}
}

// Register is used to initialize & setup the plugin
func (p *HelloWorldPlugin) Register(handler pkg.Handler) error {
	p.handler = handler
	return nil
}

func helloWorldEvent(event string, v interface{}) {
	log.Print("Hello World")
}
