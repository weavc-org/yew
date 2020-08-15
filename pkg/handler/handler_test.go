package handler

import (
	"fmt"

	"github.com/weavc/yuu/pkg/plugin"

	"testing"
)

func TestPluginLoader(t *testing.T) {
	m := NewHandler()

	test := NewTestPlugin("test1", func(m plugin.Handler) error { return nil })

	e := m.LoadPlugin(test)
	if e != nil {
		t.Error(e)
	}

	test.r = func(m plugin.Handler) error { return fmt.Errorf("an error occured") }
	e = m.LoadPlugin(test)
	if e == nil {
		t.Error(e)
	}
}

func TestGets(t *testing.T) {

	m := NewHandler()

	test1 := NewTestPlugin("test1", func(m plugin.Handler) error { return nil })
	e := m.LoadPlugin(test1)
	if e != nil {
		t.Error(e)
	}

	plgs := m.GetPlugins()
	if len(plgs) != 1 {
		t.Errorf("%d plugins loaded, expected 1", len(plgs))
	}

	test2 := NewTestPlugin("test2", func(m plugin.Handler) error { return nil })
	e = m.LoadPlugin(test2)
	if e != nil {
		t.Error(e)
	}

	servs := m.GetServices()
	if len(servs) != 2 {
		t.Errorf("%d services loaded, expected 2", len(servs))
	}

	m.Walk(func(m plugin.Manifest, v plugin.Plugin) {
		if m.Name != "test1" && m.Name != "test2" {
			t.Errorf("Incorrect plugin found in walk")
		}
	})
}

func TestEvents(t *testing.T) {
	var response string

	m := NewHandler()
	m.On("test1", func(v interface{}) {
		s := v.(string)
		response = s
		if response != "hello!" {
			t.Errorf("Incorrect response. recieved %s, expected hello!", response)
		}
	})

	test := NewTestPlugin("test1", func(m plugin.Handler) error {
		m.Emit("test1", "hello!")
		return nil
	})

	e := m.LoadPlugin(test)
	if e != nil {
		t.Error(e)
	}
}

func NewTestPlugin(name string, r func(m plugin.Handler) error) *TestPlugin1 {
	return &TestPlugin1{r: r, name: name}
}

type TestPlugin1 struct {
	r    func(m plugin.Handler) error
	name string
	plugin.Plugin
	plugin.Service
}

func (p *TestPlugin1) Manifest() plugin.Manifest {
	return plugin.Manifest{Name: p.name, Description: "Plugin used in testing"}
}

func (p *TestPlugin1) Register(m plugin.Handler) error {
	return p.r(m)
}

func (p *TestPlugin1) Start() {
	return
}

type TestPlugin2 struct{}
