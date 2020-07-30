package plugin

const (
	LOADED            Event = "initialized"
	PLUGIN_REGISTERED Event = "registered"
	SERVICE_STARTED   Event = "service"
)

type Event string

// Handler interface, handles and provides plugins with an interface to extend and communicate with the base functionality
// implemention in internal/Handler.go
type Handler interface {
	LoadPluginDir(directory string) error
	LoadPluginPath(path string) error
	LoadPlugin(v Plugin) error
	Emit(name Event, v interface{})
	On(name Event, callback func(v interface{}))
	Walk(func(manifest *Manifest, v Plugin))
	GetPlugins() []Plugin
	GetServices() []Service
}

type Manifest struct {
	Name        string
	Description string
	Registered  bool
	Events      []string
}

// Plugin defines a basic plugin
type Plugin interface {
	Manifest() *Manifest
	Register(m Handler) error
}

// Service defines a basic service plugin
// start is run in a go routine, shortly after Handler is initialized
// todo: context, channel
type Service interface {
	Plugin
	Start()
}
