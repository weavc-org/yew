package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

func FetchConfig(path string, plugin string, v interface{}) error {
	path, err := CheckPath(path)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var u map[string]interface{} = make(map[string]interface{})

	err = yaml.Unmarshal(data, &u)
	if err != nil {
		return err
	}

	if u[plugin] == nil {
		return fmt.Errorf("No config found")
	}

	err = mapstructure.Decode(u[plugin], v)
	if err != nil {
		return err
	}

	return nil
}

func CheckPath(path string) (string, error) {
	_, e := os.Stat(path)
	if e != nil {
		return "", e
	}

	o, e := filepath.Abs(path)
	if e != nil {
		return "", e
	}

	return o, nil
}
