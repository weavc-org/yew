package main

import (
	"log"

	"github.com/weavc/yuu/pkg"
)

// Plugin variable must be exported
// Plugin must also implement the github.com/weavc/yuu/pkg/plugin.Plugin interface
var Plugin HelloWorldPlugin = HelloWorldPlugin{}

type HelloWorldPlugin struct {
	handler pkg.Handler
	events  map[string]func(v interface{})

	pkg.Plugin
}

func (p *HelloWorldPlugin) Manifest() pkg.Manifest {
	return pkg.Manifest{
		Name:        "HelloWorld",
		Description: "Hello world event plugin",
		Events:      map[string]func(v interface{}){pkg.LOADED: helloWorldEvent},
	}
}

// Register is used to initialize & setup the plugin
func (p *HelloWorldPlugin) Register(m pkg.Handler) error {

	// store Handler pointer
	p.handler = m
	return nil
}

func helloWorldEvent(v interface{}) {
	log.Print("Hello World")
}
