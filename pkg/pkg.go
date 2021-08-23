package pkg

// Handler defines a set of functions that can be used by plugins and the application
// see pkg/handler for implementation of interface
type Handler interface {
	LoadPluginsDir(directory string) error
	LoadPlugins(v ...Plugin) error
	Walk(func(manifest Manifest, v Plugin))
	GetPlugins() []Plugin
}

// Manifest defines the requirements of a plugin to the handler
type Manifest struct {
	Namespace   string
	Description string
	Data        map[string]interface{}
}

// Basic plugin interface
// This contains the definition for a plugin required by the handler
// See Service for how this can be extended to provide more functionality
type Plugin interface {
	Manifest() Manifest
	Setup(h Handler) error
}

// Service defines a service plugin
// service plugins are started as a go routine when the plugin is registered
type Service interface {
	Plugin
	Start()
}
