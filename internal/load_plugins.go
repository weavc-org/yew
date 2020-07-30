package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	p "github.com/mogolade/yuu/pkg/plugin"
)

func LoadPlugins(directory string, handler func(v p.Plugin) error) error {
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

func LoadPlugin(path string, handler func(v p.Plugin) error) error {
	p1, err := plugin.Open(path)
	if err != nil {
		return err
	}

	p2, err := p1.Lookup("Plugin")
	if err != nil {
		return err
	}

	p3, t := p2.(p.Plugin)
	if t == false {
		return fmt.Errorf("Plugin variable not of correct type. should implement github.com/mogolade/yuu/pkg/plugin/Plugin")
	}

	return handler(p3)
}
