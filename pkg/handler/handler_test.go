package handler

import (
	"fmt"

	"github.com/weavc/yuu/pkg/types"

	"testing"
)

func TestPluginLoader(t *testing.T) {
	m := NewHandler(nil)

	test := NewTestPlugin("test1", func(m types.Handler) error { return nil })

	e := m.LoadPlugin(test)
	if e != nil {
		t.Error(e)
	}

	test.r = func(m types.Handler) error { return fmt.Errorf("an error occured") }
	e = m.LoadPlugin(test)
	if e == nil {
		t.Error(e)
	}
}

func TestGets(t *testing.T) {

	m := NewHandler(nil)

	test1 := NewTestPlugin("test1", func(m types.Handler) error { return nil })
	e := m.LoadPlugin(test1)
	if e != nil {
		t.Error(e)
	}

	plgs := m.GetPlugins()
	if len(plgs) != 1 {
		t.Errorf("%d plugins loaded, expected 1", len(plgs))
	}

	test2 := NewTestPlugin("test2", func(m types.Handler) error { return nil })
	e = m.LoadPlugin(test2)
	if e != nil {
		t.Error(e)
	}

	servs := m.GetServices()
	if len(servs) != 2 {
		t.Errorf("%d services loaded, expected 2", len(servs))
	}

	m.Walk(func(m types.Manifest, v types.Plugin) {
		if m.Name != "test1" && m.Name != "test2" {
			t.Errorf("Incorrect plugin found in walk")
		}
	})
}

func TestEvents(t *testing.T) {
	// var response string

	m := NewHandler(nil)
	// m.On("test1", func(v interface{}) {
	// 	s := v.(string)
	// 	response = s
	// 	if response != "hello!" {
	// 		t.Errorf("Incorrect response. recieved %s, expected hello!", response)
	// 	}
	// })

	test := NewTestPlugin("test1", func(m types.Handler) error {
		m.Emit("test1", "hello!")
		return nil
	})

	e := m.LoadPlugin(test)
	if e != nil {
		t.Error(e)
	}
}

func NewTestPlugin(name string, r func(m types.Handler) error) *TestPlugin1 {
	return &TestPlugin1{r: r, name: name}
}

type TestPlugin1 struct {
	r    func(m types.Handler) error
	name string
	types.Plugin
	types.Service
}

func (p *TestPlugin1) Manifest() types.Manifest {
	return types.Manifest{Name: p.name, Description: "Plugin used in testing"}
}

func (p *TestPlugin1) Register(m types.Handler) error {
	return p.r(m)
}

func (p *TestPlugin1) Start() {
	return
}

type TestPlugin2 struct{}
