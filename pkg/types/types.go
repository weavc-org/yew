package types

const (
	LOADED            string = "initialized"
	PLUGIN_REGISTERED string = "registered"
	SERVICE_STARTED   string = "service"
)

// Handler interface, handles and provides plugins with an interface to extend and communicate with the base functionality
// implemention in internal/Handler.go
type Handler interface {
	LoadPluginDir(directory string) error
	LoadPluginPath(path string) error
	LoadPlugin(v Plugin) error
	Emit(name string, v interface{})
	Walk(func(manifest Manifest, v Plugin))
	GetPlugins() []Plugin
	GetServices() []Service
	LoadConfig(plg Plugin, v interface{}) error
	SetupConfigDirectory(dir string) error
}

type Manifest struct {
	Name         string
	Description  string
	Events       map[string]func(v interface{})
	Config       interface{}
	CreateConfig func() interface{}
	Data         map[string]interface{}
	Registered   bool
}

// Plugin defines a basic plugin
type Plugin interface {
	Manifest() Manifest
	Register(m Handler) error
}

// Service defines a basic service plugin
// start is run in a go routine, shortly after Handler is initialized
// todo: context, channel
type Service interface {
	Plugin
	Start()
}
