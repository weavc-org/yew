package main

import (
	"log"

	"github.com/weavc/yuu/pkg/plugin"
)

// Plugin variable must be exported
// Plugin must also implement the github.com/weavc/yuu/pkg/plugin.Plugin interface
var Plugin HelloWorldPlugin = HelloWorldPlugin{}

type HelloWorldPlugin struct {
	handler plugin.Handler

	plugin.Plugin
}

func (p *HelloWorldPlugin) Manifest() plugin.Manifest {
	return plugin.Manifest{Name: "HelloWorld", Description: "Hello world event plugin"}
}

// Register is used to initialize & setup the plugin
func (p *HelloWorldPlugin) Register(m plugin.Handler) error {

	// store Handler pointer
	p.handler = m

	// register event
	p.handler.On(plugin.LOADED, helloWorldEvent)

	return nil
}

func helloWorldEvent(v interface{}) {
	log.Print("Hello World")
}
