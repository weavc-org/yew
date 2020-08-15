package handler

import (
	"fmt"

	"github.com/weavc/yew/internal"
	"github.com/weavc/yew/internal/configs"
	"github.com/weavc/yew/pkg"
)

// Handler struct. This implements plugin/Handler interface and handles the plugins
type Handler struct {
	pkg.Handler

	plugins []pkg.Plugin

	Config *Config
}

// LoadPluginsDir will load all plugins in provided directory path
func (m *Handler) LoadPluginsDir(directory string) error {
	e := internal.LoadPlugins(directory, m.loadPlugin)
	if e != nil {
		return e
	}
	m.Emit(pkg.LOADED, nil)
	return nil
}

func (m *Handler) LoadPlugins(p ...pkg.Plugin) error {
	for _, plugin := range p {
		err := m.loadPlugin(plugin)
		if err != nil {
			return fmt.Errorf("Error loading plugin: %s", err)
		}
	}

	return nil
}

// Walk will take you on a walk through the registered plugins
// for each plugin, the handler function passed through will be called
func (m *Handler) Walk(handler func(manifest pkg.Manifest, v pkg.Plugin)) {
	for _, p := range m.plugins {
		handler(p.Manifest(), p)
	}
}

// GetPlugins will return an array of registered plugins
func (m *Handler) GetPlugins() []pkg.Plugin {
	var plgs []pkg.Plugin
	m.Walk(func(manifest pkg.Manifest, v pkg.Plugin) {
		p, t := v.(pkg.Plugin)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// NewHandler creates & returns pkg.Handler structure
func NewHandler(c *Config) pkg.Handler {
	if c == nil {
		c = DefaultConfig
	}

	if c.PluginConfigDirectory != "" {
		s, err := configs.CheckPath(c.PluginConfigDirectory)
		if err != nil {
			panic(err)
		}

		c.PluginConfigDirectory = s
	}

	m := &Handler{Config: c}
	return m
}

// LoadPlugin will take a struct that implements plugin.Plugin and load it into the Handler
func (m *Handler) loadPlugin(p pkg.Plugin) error {
	plg, c := p.(pkg.Plugin)
	if c == true {
		m.plugins = append(m.plugins, plg)

		man := plg.Manifest()

		e := plg.Register(m)
		if e != nil {
			return e
		}

		m.Emit(pkg.PLUGIN_REGISTERED, plg)

		if man.Config != nil {
			err := m.LoadConfig(plg, man.Config)
			if err != nil {
				fmt.Print(fmt.Errorf("there was an error loading config for %s. this could mean the file was missing, or there were errors loading it", man.Name))
			}
		}

		service, c := p.(pkg.Service)
		if c == true {
			if m.Config.Services == true {
				go service.Start()
				m.Emit(pkg.SERVICE_STARTED, plg)
			}
		}

		return nil
	}

	return fmt.Errorf("Plugin does not implement Plugin interface")
}
