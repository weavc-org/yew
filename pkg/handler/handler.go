package handler

import (
	"fmt"

	"github.com/weavc/yuu/internal"
	"github.com/weavc/yuu/internal/configs"
	"github.com/weavc/yuu/pkg/types"
)

// Handler struct. This implements plugin/Handler interface and handles the plugins
type Handler struct {
	types.Handler

	plugins []types.Plugin

	Config *HandlerConfig
}

// LoadPluginPath will load a plugin via a file path
func (m *Handler) LoadPluginPath(path string) error {
	e := internal.LoadPlugin(path, m.LoadPlugin)
	if e != nil {
		return e
	}
	m.Emit(types.LOADED, nil)
	return nil
}

// LoadPluginDir will load all plugins in provided directory path
func (m *Handler) LoadPluginDir(directory string) error {
	e := internal.LoadPlugins(directory, m.LoadPlugin)
	if e != nil {
		return e
	}
	m.Emit(types.LOADED, nil)
	return nil
}

// LoadPlugin will take a struct that implements plugin.Plugin and load it into the Handler
func (m *Handler) LoadPlugin(v types.Plugin) error {
	plg, c := v.(types.Plugin)
	if c == true {
		m.plugins = append(m.plugins, plg)

		man := plg.Manifest()
		if man.Registered == false {
			e := plg.Register(m)
			if e != nil {
				return e
			}

			m.Emit(types.PLUGIN_REGISTERED, plg)
		}

		if man.Config != nil {
			err := m.LoadConfig(plg, man.Config)
			if err != nil {

				panic(err)
			}
		}

		service, c := v.(types.Service)
		if c == true {
			if m.Config.Services == true {
				go service.Start()
				m.Emit(types.SERVICE_STARTED, plg)
			}
		}

		man.Registered = true
		return nil
	}

	return fmt.Errorf("Plugin does not implement Plugin interface")
}

// Walk will take you on a walk through the registered plugins
// for each plugin, the handler function passed through will be called
func (m *Handler) Walk(handler func(manifest types.Manifest, v types.Plugin)) {
	for _, p := range m.plugins {
		handler(p.Manifest(), p)
	}
}

// GetPlugins will return an array of registered plugins
func (m *Handler) GetPlugins() []types.Plugin {
	var plgs []types.Plugin
	m.Walk(func(manifest types.Manifest, v types.Plugin) {
		p, t := v.(types.Plugin)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// GetServices will return an array of Service plugins
func (m *Handler) GetServices() []types.Service {
	var plgs []types.Service
	m.Walk(func(manifest types.Manifest, v types.Plugin) {
		p, t := v.(types.Service)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// NewHandler creates & returns pkg.Handler structure
func NewHandler(c *HandlerConfig) types.Handler {
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

// BuildPlugin builds package into a plugin
func BuildPlugin(output string, dir string) error {
	return internal.BuildPlugin(output, dir)
}

// BuildPlugins accepts multiple directories to be built
func BuildPlugins(output string, dirs ...string) error {
	return internal.BuildPlugins(output, dirs)
}
