package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/weavc/yew/pkg"
)

// LoadPlugins finds all plugins in the provided directory
// and then passes them through to LoadPlugin
func LoadPlugins(directory string, handler func(v pkg.Plugin) error) error {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".so" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, f := range files {
		err = LoadPlugin(f, handler)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPlugin opens the .so file, perfroms a lookup for the exported 'Plugin variable'
// and checks it typing before passing it on the the handler
func LoadPlugin(path string, handler func(v pkg.Plugin) error) error {
	p1, err := plugin.Open(path)
	if err != nil {
		return err
	}

	p2, err := p1.Lookup("Plugin")
	if err != nil {
		return err
	}

	p3, t := p2.(pkg.Plugin)
	if t == false {
		return fmt.Errorf("Plugin variable not of correct type. should implement github.com/weavc/yew/pkg/plugin/Plugin")
	}

	return handler(p3)
}
