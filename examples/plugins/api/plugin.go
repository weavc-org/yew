package main

import (
	"encoding/json"
	"net/http"

	"github.com/weavc/yuu/pkg/plugin"
)

// Plugin variable that implements github.com/weavc/yuu/pkg/plugin.Plugin
// must be exported if building into a .so file.
//This is how the Plugin is found within the binary plugin
var Plugin ApiPlugin = ApiPlugin{}

// ApiPlugin is the struct that implements plugin.Plugin & more
type ApiPlugin struct {
	handler plugin.Handler

	plugin.Plugin
}

// Manifest gives the handler & other plugins an idea of what this plugin is
func (p *ApiPlugin) Manifest() plugin.Manifest {
	return plugin.Manifest{Name: "Api", Description: "Api plugin", Events: []string{"api"}}
}

// Register is used to initialize & setup the plugin
func (p *ApiPlugin) Register(m plugin.Handler) error {
	// store Handler pointer
	p.handler = m

	return nil
}

// RegisterRoutes implements an interface defined in examples/main.go
// An example of how plugins can be extended to provide additional
// communication with different applications.
func (p ApiPlugin) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		p.handler.Emit("api", r.URL.String())
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})
}
