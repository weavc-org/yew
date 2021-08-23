package main

import (
	"log"
	"net/http"

	"github.com/weavc/yew/v2/pkg"
	"github.com/weavc/yew/v2/pkg/builders"
	"github.com/weavc/yew/v2/pkg/handler"
)

type RegisterAPI interface {
	RegisterRoutes(mux *http.ServeMux)
}

func main() {

	var mux *http.ServeMux = http.DefaultServeMux

	// build plugins from source to examples/.bin
	builders.BuildPlugins("examples/.bin/", "examples/plugins/hello-world", "examples/plugins/api")

	// create/get new Handler structure & load plugins from examples/.bin
	m := handler.NewHandler()

	m.LoadPluginsDir("examples/.bin/")

	// recieve api events, these are emitted in the api plugin
	// walk through plugins (foreach registered plugin)
	m.Walk(func(man pkg.Manifest, plugin pkg.Plugin) {
		// check if plugin implements RegisterAPI interface defined above
		p, e := plugin.(RegisterAPI)
		if e {
			// let plugin register handlers
			p.RegisterRoutes(mux)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
