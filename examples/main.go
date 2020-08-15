package main

import (
	"log"
	"net/http"

	"github.com/weavc/yuu/pkg"
	"github.com/weavc/yuu/pkg/builders"
	"github.com/weavc/yuu/pkg/handler"
)

type RegisterAPI interface {
	RegisterRoutes(mux *http.ServeMux)
}

func main() {

	var s *http.ServeMux = http.DefaultServeMux

	// build plugins from source to examples/.bin
	builders.BuildPlugins("examples/.bin/", "examples/plugins/hello-world", "examples/plugins/api")

	// create/get new Handler structure & load plugins from examples/.bin
	m := handler.NewHandler(&handler.Config{
		Services:              true,
		PluginConfigDirectory: "examples/configs",
		Events:                map[string]func(v interface{}){"api": apiEvent},
	})

	m.LoadPluginsDir("examples/.bin/")

	// recieve api events, these are emitted in the api plugin
	// walk through plugins (foreach registered plugin)
	m.Walk(func(man pkg.Manifest, plgin pkg.Plugin) {
		// check if plugin implements RegisterAPI interface defined above
		p, e := plgin.(RegisterAPI)
		if e == true {
			// let plugin register handlers
			p.RegisterRoutes(s)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiEvent(v interface{}) {
	s := v.(string)
	log.Print(s)
}
