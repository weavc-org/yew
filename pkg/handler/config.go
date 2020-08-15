package handler

import (
	"fmt"

	"github.com/weavc/yuu/internal/configs"
	"github.com/weavc/yuu/pkg/plugin"
)

type handlerConfig struct {
	name  string
	types int32
}

func (m *Handler) SetupConfigDirectory(dir string) error {
	path, err := configs.CheckPath(dir)
	if err != nil {
		return fmt.Errorf("Failed to get directory: %s", err)
	}

	m.configDir = path

	// try to load handler config
	c := defaultConfig()
	err = configs.LoadConfig(m.configDir, "yuu", c)
	// ignore errors, set config if its there & readable
	if err == nil {
		m.config = c
	}

	return nil
}

func (m *Handler) LoadConfig(plg plugin.Plugin, v interface{}) error {
	if m.configDir == "" {
		return fmt.Errorf("No config directory set")
	}

	err := configs.LoadConfig(m.configDir, plg.Manifest().Name, v)
	if err != nil {
		return fmt.Errorf("Error loading in config: %s", err)
	}

	return nil
}

func defaultConfig() *handlerConfig {
	return &handlerConfig{
		name:  "hello",
		types: 122333,
	}
}
