package pkg

const (
	// LoadedEvent event name
	// Triggered after loading a plugin or set of plugins
	LoadedEvent string = "loaded_event"
	// PluginRegisteredEvent event name
	// Triggered after registering a plugin
	PluginRegisteredEvent string = "plugin_registered_event"
	// ServiceStartedEvent event name
	// Triggered after starting a service
	ServiceStartedEvent string = "service_started_event"
)

// Handler defines a set of functions that can be used by plugins and the application
// see pkg/handler for implementation of interface
type Handler interface {
	LoadPluginsDir(directory string) error
	LoadPlugins(v ...Plugin) error
	Emit(name string, v interface{})
	Walk(func(manifest Manifest, v Plugin))
	GetPlugins() []Plugin
}

// Manifest defines the requirements of a plugin to the handler
type Manifest struct {
	Namespace   string
	Description string
	Events      map[string]func(event string, v interface{})
	Config      interface{}
	Data        map[string]interface{}
}

// Plugin defines a basic plugin
type Plugin interface {
	Manifest() Manifest
	Register(h Handler) error
}

// Service defines a service plugin
// service plugins are started as a go routine when the plugin is registered
type Service interface {
	Plugin
	Start()
}
