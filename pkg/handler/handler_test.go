package handler

import (
	"fmt"

	"testing"

	"github.com/weavc/yew/v3/pkg"
)

func TestPluginLoader(t *testing.T) {
	m := NewHandler()

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

	m := NewHandler()

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
		if m.Namespace != "test1" && m.Namespace != "test2" {
			t.Errorf("Incorrect plugin found in walk")
		}
	})
}

func TestDuplicateNamespaces(t *testing.T) {
	m := NewHandler()

	test1 := NewTestPlugin("test1", func(m pkg.Handler) error { return nil })
	test2 := NewTestPlugin("test1", func(m pkg.Handler) error { return nil })
	e := m.LoadPlugins(test1, test2)
	if e == nil {
		t.Errorf("No error when namespaces collide")
	}
}

func NewTestPlugin(name string, r func(m pkg.Handler) error) *MockPlugin1 {
	return &MockPlugin1{r: r, name: name}
}

type MockPlugin1 struct {
	r    func(m pkg.Handler) error
	name string
	pkg.Plugin
	pkg.Service
}

func (p *MockPlugin1) Manifest() pkg.Manifest {
	return pkg.Manifest{Namespace: p.name, Description: "Plugin used in testing"}
}

func (p *MockPlugin1) Setup(m pkg.Handler) error {
	return p.r(m)
}

func (p *MockPlugin1) Start() {
	return
}

type MockPlugin2 struct{}
