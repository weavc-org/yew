package main

import (
	"encoding/json"
	"net/http"

	"github.com/weavc/yew/v2/pkg"
)

// Plugin variable that implements github.com/weavc/yew/v2/pkg/plugin.Plugin
// must be exported if building into a .so file.
//This is how the Plugin is found within the binary plugin
var Plugin APIPlugin = APIPlugin{config: &c{}}

// APIPlugin is the struct that implements plugin.Plugin & more
type APIPlugin struct {
	handler pkg.Handler

	pkg.Plugin
	config *c
}

// Manifest gives the handler & other plugins an idea of what this plugin is
func (p *APIPlugin) Manifest() pkg.Manifest {
	return pkg.Manifest{Namespace: "examples.api", Description: "Api plugin", Config: p.config}
}

// Register is used to initialize & setup the plugin
func (p *APIPlugin) Register(m pkg.Handler) error {
	// store Handler pointer
	p.handler = m

	return nil
}

// RegisterRoutes implements an interface defined in examples/main.go
// An example of how plugins can be extended to provide additional
// communication with different applications.
func (p *APIPlugin) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		p.handler.Emit("api", r.URL.String())
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "payload": p.config.Payload})
	})
}

type c struct {
	Payload string
}
