![tests](https://github.com/weavc/yew/workflows/Go/badge.svg?branch=master) 
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/weavc/yew)

#### Yew

Yew is a lightweight module mostly in its development/testing phase. Its aim is to aid in using plugin/event driven architecture within go. There are probably alternatives out there, I just started doing something small and it turned into this.

There are 2 main parts to this module, the `Plugin` and `Handler`. The `Handler` can register any `Plugin`'s it is asked to load, this can be through `.so` file(s) (built using `go build -buildmode=plugin`) or any struct that implements the `pkg/plugin.Plugin` interface, example of this [here](#Plugins).

#### Handler

The Handler is responsible for handling plugins and providing a number of utiliies to both the plugins and the application implementing the handler. Plugins can be loaded via any stuct that implements the `pkg.Plugin` interface, or `.so` files where they have an export `Plugin` variable the implements `pkg.Plugin`. Both the implementing application and the plugins will have access to the handler, but only the implementing application should import the `pkg/handler` package. The handler will pass itself to the plugins when they are registered via the `Register` function on the Plugins interface, this allows them to store the reference and use it for communicating. 

Example use of the handler: 

```go
package main

import (
	"log"

    "github.com/weavc/yew/pkg"
    "github.com/weavc/yew/pkg/handler"
)

var h pkg.Handler

func main() {
	// create/get new Handler structure & load plugins from examples/.bin
	h = handler.NewHandler(&handler.Config{
		Services: true,
		Events:   map[string]func(event string, v interface{}){pkg.LoadedEvent: onLoad},
	})

	// load .so files inside '.plugins/' dir
	h.LoadPluginsDir(".plugins/")

	h.LoadPlugins(&Plugin)

	t := time.Second * 5
	<-time.After(t)
}

func onLoad(event string, v interface{}) {
	h.Walk(func(manifest pkg.Manifest, plgin pkg.Plugin) {
		log.Printf("Loaded: %s", manifest.Namespace)
	})
}
```

#### Plugins

Plugins should only ever import `github.com/weavc/yew/pkg`, this helps reduce circular reference issues and also the need to rebuild for any minor releases. If the plugin is being built & distibuted via `.so` files (built using `go build -buildmode=plugin`), there should be an exported variable named `Plugin` in the main package, this is how the handler will find the Plugin in the binary, see below example for what this should look like.

Example plugin:
```go
// this variable is looked up when loading the plugin from a .so file
// so make sure its there!
var Plugin plugin = plugin{}

// simple plugin struct, implements pkg.Plugin
type plugin struct {
	pkg.Plugin
	handler pkg.Handler
}

// Called as its being registered to the handler, can be used for setup/initialization
// as this is the first thing that happens.
// Also recommended to store the pkg.Handler reference at this stage
func (p *plugin) Register(handler pkg.Handler) error {
	p.handler = handler
	return nil
}

// returns the manifest which defines the plugin, what events it would like to recieve, config etc
func (p *plugin) Manifest() pkg.Manifest {
	return pkg.Manifest{Namespace: "plugin.example", Description: "Example plugin"}
}
```

Plugins can also be loaded directly via their structure using `Handler.LoadPlugins(<&plugin{}>)`.