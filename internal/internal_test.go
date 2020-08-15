package internal

import (
	"testing"

	"github.com/weavc/yuu/pkg"
)

var built bool = false

func TestBuildPlugins(t *testing.T) {
	e := buildPlugins()
	if e != nil {
		t.Error(e)
	}
}

func TestBuildPluginsFakeOutput(t *testing.T) {
	e := BuildPlugins("./0da4d1b6e798c53b07b1d71a13773822", []string{"../examples/plugins/api", "../examples/plugins/hello-world"})
	if e == nil {
		t.Error("expected error")
	}
}

func TestBuildPluginsFakePluginDir(t *testing.T) {
	e := BuildPlugins("../examples/.bin/", []string{"./0da4d1b6e798c53b07b1d71a13773822"})
	if e == nil {
		t.Error("expected error")
	}
}

func TestLoadPlugins(t *testing.T) {
	if built == false {
		e := buildPlugins()
		if e != nil {
			t.Error(e)
		}
	}

	e := LoadPlugins("../examples/.bin/", func(v pkg.Plugin) error { return nil })
	if e != nil {
		t.Error(e)
	}
}

func buildPlugins() error {
	e := BuildPlugins("../examples/.bin/", []string{"../examples/plugins/api", "../examples/plugins/hello-world"})
	built = true
	return e
}
