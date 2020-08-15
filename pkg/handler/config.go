package handler

import (
	"fmt"

	"github.com/weavc/yuu/internal/configs"
	"github.com/weavc/yuu/pkg/types"
)

var DefaultConfig *HandlerConfig = &HandlerConfig{Services: true}

type HandlerConfig struct {
	Services              bool
	PluginConfigDirectory string
	Events                map[string]func(v interface{})
}

func (m *Handler) LoadConfig(plg types.Plugin, v interface{}) error {
	if m.Config.PluginConfigDirectory == "" {
		return fmt.Errorf("No config directory set")
	}

	err := configs.LoadConfig(m.Config.PluginConfigDirectory, plg.Manifest().Name, v)
	if err != nil {
		return fmt.Errorf("Error loading in config: %s", err)
	}

	fmt.Print(v)

	return nil
}
