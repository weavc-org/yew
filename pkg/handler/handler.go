package handler

import (
	"fmt"
	"log"

	"github.com/weavc/yew/v2/internal"
	"github.com/weavc/yew/v2/internal/configs"
	"github.com/weavc/yew/v2/pkg"
)

// Handler implements pkg.Handler
type Handler struct {
	pkg.Handler
	plugins []pkg.Plugin
	Config  *Config
}

// LoadPluginsDir will load all .so plugins in given directory
func (h *Handler) LoadPluginsDir(directory string) error {
	err := internal.LoadPlugins(directory, h.loadPlugin)
	if err != nil {
		return err
	}
	h.Emit(pkg.LoadedEvent, nil)
	return nil
}

// LoadPlugins takes plugins, loads and registers them with the handler
func (h *Handler) LoadPlugins(plugins ...pkg.Plugin) error {
	for _, plugin := range plugins {
		err := h.loadPlugin(plugin)
		if err != nil {
			return fmt.Errorf("Error loading plugin: %s", err)
		}
	}

	return nil
}

// Walk will take you on a walk through the registered plugins
// for each plugin, the handler function passed through will be called
func (h *Handler) Walk(handler func(manifest pkg.Manifest, plugin pkg.Plugin)) {
	for _, p := range h.plugins {
		handler(p.Manifest(), p)
	}
}

// GetPlugins will return an array of registered plugins
func (h *Handler) GetPlugins() []pkg.Plugin {
	var plgs []pkg.Plugin
	h.Walk(func(manifest pkg.Manifest, v pkg.Plugin) {
		p, t := v.(pkg.Plugin)
		if t == true {
			plgs = append(plgs, p)
		}
	})

	return plgs
}

// NewHandler creates & returns a new pkg.Handler struct
func NewHandler(c *Config) pkg.Handler {
	if c == nil {
		c = DefaultConfig
	}

	if c.PluginConfigPath != "" {
		s, err := configs.CheckPath(c.PluginConfigPath)
		if err != nil {
			panic(err)
		}

		c.PluginConfigPath = s
	}

	m := &Handler{Config: c}
	return m
}

// loadPlugin is the initialization handler for plugins, triggered after the plugin is loaded
func (h *Handler) loadPlugin(p pkg.Plugin) error {
	plg, c := p.(pkg.Plugin)
	if c == false {
		return fmt.Errorf("Plugin does not implement the pkg.Plugin interface")
	}

	man := plg.Manifest()

	if h.Config.UniqueNamespaces {
		for _, fplg := range h.plugins {
			if fplg.Manifest().Namespace == man.Namespace {
				return fmt.Errorf("Plugin namespace clash %s", man.Namespace)
			}
		}
	}

	e := plg.Register(h)
	if e != nil {
		return e
	}

	h.plugins = append(h.plugins, plg)
	h.Emit(pkg.PluginRegisteredEvent, plg)

	if man.Config != nil {
		err := h.FetchConfig(plg, man.Config)
		if err != nil {
			log.Printf("There was an error loading config for %s: %v", man.Namespace, err)
		}
	}

	service, c := p.(pkg.Service)
	if c == true {
		if h.Config.Services == true {
			go service.Start()
			h.Emit(pkg.ServiceStartedEvent, plg)
		}
	}

	return nil
}
