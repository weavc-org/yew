package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig(dir string, name string, v interface{}) error {
	path, err := CheckPath(dir)
	if err != nil {
		return err
	}

	vip := viper.New()
	vip.SetConfigName(name)

	// this will come from handler
	vip.AddConfigPath(path)

	err = vip.ReadInConfig()
	if err != nil {
		return err
	}

	err = vip.Unmarshal(v)
	if err != nil {
		return err
	}

	return nil
}

func CheckPath(dir string) (string, error) {
	info, e := os.Stat(dir)
	if e != nil {
		return "", e
	}

	if !info.IsDir() {
		return "", fmt.Errorf("output is not a directory")
	}

	o, e := filepath.Abs(dir)
	if e != nil {
		return "", e
	}

	return o, nil
}
