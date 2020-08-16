package main

import (
	"log"

	"github.com/weavc/yew/v2/pkg"
)

// Plugin variable must be exported
// Plugin must also implement the github.com/weavc/yew/v2/pkg/plugin.Plugin interface
var Plugin HelloWorldPlugin = HelloWorldPlugin{}

type HelloWorldPlugin struct {
	handler pkg.Handler
	pkg.Plugin

	c *Config
}

type Config struct {
	Say string
}

// Manifest returns
func (p *HelloWorldPlugin) Manifest() pkg.Manifest {
	return pkg.Manifest{
		Namespace:   "examples.helloworld",
		Description: "Hello world event plugin",
		Config:      &p.c,
		Events:      map[string]func(event string, v interface{}){pkg.LoadedEvent: p.helloWorldEvent},
	}
}

// Register is used to initialize & setup the plugin
func (p *HelloWorldPlugin) Register(handler pkg.Handler) error {
	p.c = &Config{Say: "earth"}
	p.handler = handler
	return nil
}

func (p *HelloWorldPlugin) helloWorldEvent(event string, v interface{}) {
	log.Printf("Hello %s", p.c.Say)
}
