package handler

import (
	"fmt"

	"testing"

	"github.com/weavc/yew/pkg"
)

func TestPluginLoader(t *testing.T) {
	m := NewHandler(nil)

	test := NewTestPlugin("test1", func(m pkg.Handler) error { return nil })

	e := m.LoadPlugins(test)
	if e != nil {
		t.Error(e)
	}

	test.r = func(m pkg.Handler) error { return fmt.Errorf("an error occured") }
	e = m.LoadPlugins(test)
	if e == nil {
		t.Error(e)
	}
}

func TestGets(t *testing.T) {

	m := NewHandler(nil)

	test1 := NewTestPlugin("test1", func(m pkg.Handler) error { return nil })
	e := m.LoadPlugins(test1)
	if e != nil {
		t.Error(e)
	}

	plgs := m.GetPlugins()
	if len(plgs) != 1 {
		t.Errorf("%d plugins loaded, expected 1", len(plgs))
	}

	test2 := NewTestPlugin("test2", func(m pkg.Handler) error { return nil })
	e = m.LoadPlugins(test2)
	if e != nil {
		t.Error(e)
	}

	m.Walk(func(m pkg.Manifest, v pkg.Plugin) {
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

	test := NewTestPlugin("test1", func(m pkg.Handler) error {
		m.Emit("test1", "hello!")
		return nil
	})

	e := m.LoadPlugins(test)
	if e != nil {
		t.Error(e)
	}
}

func NewTestPlugin(name string, r func(m pkg.Handler) error) *TestPlugin1 {
	return &TestPlugin1{r: r, name: name}
}

type TestPlugin1 struct {
	r    func(m pkg.Handler) error
	name string
	pkg.Plugin
	pkg.Service
}

func (p *TestPlugin1) Manifest() pkg.Manifest {
	return pkg.Manifest{Name: p.name, Description: "Plugin used in testing"}
}

func (p *TestPlugin1) Register(m pkg.Handler) error {
	return p.r(m)
}

func (p *TestPlugin1) Start() {
	return
}

type TestPlugin2 struct{}
