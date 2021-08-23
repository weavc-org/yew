package handler

import (
	"fmt"

	"github.com/weavc/yew/v2/internal"
	"github.com/weavc/yew/v2/pkg"
)

// Handler implements pkg.Handler
type Handler struct {
	pkg.Handler
	plugins map[string]pkg.Plugin
}

// LoadPluginsDir will load all .so plugins in given directory
func (h *Handler) LoadPluginsDir(directory string) error {
	err := internal.LoadPlugins(directory, h.loadPlugin)
	if err != nil {
		return err
	}
	return nil
}

// LoadPlugins takes plugins, loads and registers them with the handler
func (h *Handler) LoadPlugins(plugins ...pkg.Plugin) error {
	for _, plugin := range plugins {
		err := h.loadPlugin(plugin)
		if err != nil {
			return fmt.Errorf("error loading plugin: %s", err)
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
		plgs = append(plgs, v)
	})

	return plgs
}

// NewHandler creates & returns a new pkg.Handler struct
func NewHandler() pkg.Handler {
	m := &Handler{plugins: make(map[string]pkg.Plugin)}
	return m
}

// loadPlugin is the initialization handler for plugins, triggered after the plugin is loaded
func (h *Handler) loadPlugin(p pkg.Plugin) error {

	_, exists := h.plugins[p.Manifest().Namespace]
	if exists {
		return fmt.Errorf("plugin namespace clash %s", p.Manifest().Namespace)
	}

	e := p.Setup(h)
	if e != nil {
		return e
	}

	h.plugins[p.Manifest().Namespace] = p

	service, c := p.(pkg.Service)
	if c {
		go service.Start()
	}

	return nil
}
