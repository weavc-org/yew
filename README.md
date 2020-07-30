![tests](https://github.com/mogolade/yuu/workflows/Go/badge.svg?branch=master) 
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/mogolade/yuu)

#### Yuu

Yuu is a lightweight module in its early stages of development/testing. Its aim is to aid in using plugin/event driven architecture within go. There are probably alternatives out there, I just started doing something small and it turned into this.

There are 2 main parts to this module, the `Plugin` and `Handler`. The `Handler` can register any `Plugin`'s it is asked to load, this can be through `.so` file(s) (built using `go build -buildmode=plugin`) or any struct that implements the `pkg/plugin.Plugin` interface, example of this [here](#Plugins).

#### Handler

The Handler is used to manage any loaded plugins. Plugins can be loaded via any stuct that implements the `pkg/plugin.Plugin` interface, or `.so` files. This should be used in any main application and not really in the plugins themselves. The handler will pass itself to the plugins when they are registered via the `Register` function on the Plugins interface, this allows them to store the managers pointer reference and use it for communicating with other plugins.

Example use of the handler: 

```go
import (
    "github.com/mogolade/yuu/pkg"
    "github.com/mogolade/yuu/pkg/plugin"
)

type RegisterAPI interface {
	RegisterRoutes(mux *http.ServeMux)
}

func main() {

    var s *http.ServeMux = http.DefaultServeMux

    // build plugins from source to examples/.bin
	pkg.BuildPlugins("examples/.bin/", []string{"examples/plugins/hello-world", "examples/plugins/api"})

	// create/get new Handler structure & load plugins from examples/.bin
	m := pkg.NewHandler()
	m.LoadPluginDir("examples/.bin/")

    // recieve api events, these are emitted in the api plugin
	m.On("api", func(v interface{}) {
		s := v.(string)
		log.Print(s)
	})

	// walk through plugins (foreach registered plugin)
	m.Walk(func(man *plugin.Manifest, plgin plugin.Plugin) {
		// check if plugin implements RegisterAPI interface defined above
		p, e := plgin.(RegisterAPI)
		if e == true {
			// let plugin register handlers
			p.RegisterRoutes(s)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

#### Plugins

Plugins should only ever import `"github.com/mogolade/yuu/pkg/plugin"`, this helps reduce circular reference issues and also the need to rebuild for any minor releases. If the plugin is being built & distibuted via the `.so` file (built using `go build -buildmode=plugin`), there should be an exported variable named `Plugin` in the main package, this is how the handler will find the Plugin in the binary, see below example for what this should look like.

Example plugin:

```go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/mogolade/yuu/pkg/plugin"
)

// Plugin variable that implements github.com/mogolade/yuu/pkg/plugin.Plugin
// must be exported if building into a .so file.
//This is how the Plugin is found within the binary plugin
var Plugin ApiPlugin = ApiPlugin{}

// ApiPlugin is the struct that implements plugin.Plugin & more
type ApiPlugin struct {
	handler plugin.Handler

	plugin.Plugin
}

// Manifest gives the handler & other plugins an idea of what this plugin is
func (p *ApiPlugin) Manifest() *plugin.Manifest {
	return &plugin.Manifest{Name: "Api", Description: "Api plugin", Events: []string{"api"}}
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
```
