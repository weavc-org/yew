package handler

import (
	"fmt"

	"github.com/weavc/yew/v2/internal/configs"
	"github.com/weavc/yew/v2/pkg"
)

// DefaultConfig is the default configuration for the handler
// this is used if nil is sent through to handler.NewHandler
var DefaultConfig *Config = &Config{Services: true, UniqueNamespaces: true}

// Config struct, defines what the handler does
type Config struct {
	Services         bool
	UniqueNamespaces bool
	PluginConfigPath string
	Events           map[string]func(event string, v interface{})
}

// FetchConfig will get the config from the file and bind it to config
func (h *Handler) FetchConfig(plg pkg.Plugin, config interface{}) error {
	if h.Config.PluginConfigPath == "" {
		return fmt.Errorf("No config set")
	}

	err := configs.FetchConfig(h.Config.PluginConfigPath, plg.Manifest().Namespace, config)
	if err != nil {
		return fmt.Errorf("Error loading in config: %s", err)
	}

	return nil
}
