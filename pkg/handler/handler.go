package handler

import (
	"fmt"

	"github.com/weavc/yuu/internal"
	"github.com/weavc/yuu/pkg/plugin"
)

// Handler struct. This implements plugin/Handler interface and handles the plugins
type Handler struct {
	plugin.Handler

	eventHandlers []*eventHandler

	plugins []plugin.Plugin

	config    *handlerConfig
	configDir string
}

// LoadPluginPath will load a plugin via a file path
func (m *Handler) LoadPluginPath(path string) error {
	e := internal.LoadPlugin(path, m.LoadPlugin)
	if e != nil {
		return e
	}
	m.Emit(plugin.LOADED, nil)
	return nil
}

// LoadPluginDir will load all plugins in provided directory path
func (m *Handler) LoadPluginDir(directory string) error {
	e := internal.LoadPlugins(directory, m.LoadPlugin)
	if e != nil {
		return e
	}
	m.Emit(plugin.LOADED, nil)
	return nil
}

// LoadPlugin will take a struct that implements plugin.Plugin and load it into the Handler
func (m *Handler) LoadPlugin(v plugin.Plugin) error {
	plg, c := v.(plugin.Plugin)
	if c == true {
		m.plugins = append(m.plugins, plg)

		man := plg.Manifest()
		if man.Registered == false {
			e := plg.Register(m)
			if e != nil {
				return e
			}

			m.Emit(plugin.PLUGIN_REGISTERED, plg)
		}

		service, c := v.(plugin.Service)
		if c == true {
			go service.Start()
			m.Emit(plugin.SERVICE_STARTED, plg)
		}

		man.Registered = true
		return nil
	}

	return fmt.Errorf("Plugin does not implement Plugin interface")
}

// On registers a callback function, triggered on each emit of an event with the given name
func (m *Handler) On(name plugin.Event, callback func(v interface{})) {
	m.eventHandlers = append(m.eventHandlers, &eventHandler{Name: name, Callback: callback})
}

// Emit emits an event, triggering any registered callbacks
// for events of the same name registered through On
func (m *Handler) Emit(name plugin.Event, v interface{}) {
	for _, h := range m.eventHandlers {
		if h.Name == name {
			if h.Callback != nil {
				go h.Callback(v)
			}
		}
	}
}

// Walk will take you on a walk through the registered plugins
// for each plugin, the handler function passed through will be called
func (m *Handler) Walk(handler func(manifest plugin.Manifest, v plugin.Plugin)) {
	for _, p := range m.plugins {
		handler(p.Manifest(), p)
	}
}

// GetPlugins will return an array of registered plugins
func (m *Handler) GetPlugins() []plugin.Plugin {
	var plgs []plugin.Plugin
	m.Walk(func(manifest plugin.Manifest, v plugin.Plugin) {
		p, t := v.(plugin.Plugin)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// GetServices will return an array of Service plugins
func (m *Handler) GetServices() []plugin.Service {
	var plgs []plugin.Service
	m.Walk(func(manifest plugin.Manifest, v plugin.Plugin) {
		p, t := v.(plugin.Service)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// NewHandler creates & returns pkg.Handler structure
func NewHandler() plugin.Handler {
	m := &Handler{config: defaultConfig()}
	return m
}

// BuildPlugin builds package into a plugin
func BuildPlugin(output string, dir string) error {
	return internal.BuildPlugin(output, dir)
}

// BuildPlugins accepts multiple directories to be built
func BuildPlugins(output string, dirs []string) error {
	return internal.BuildPlugins(output, dirs)
}

type eventHandler struct {
	Name     plugin.Event
	Callback func(v interface{})
}
