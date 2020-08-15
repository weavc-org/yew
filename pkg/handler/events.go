package handler

import (
	"github.com/weavc/yuu/pkg/types"
)

// Emit emits an event, triggering any registered callbacks
// for events of the same name
func (m *Handler) Emit(name string, v interface{}) {
	m.Walk(func(m types.Manifest, v types.Plugin) {
		if m.Events != nil {
			for key, handler := range m.Events {
				if key == name {
					go handler(v)
				}
			}
		}
	})

	for key, handler := range m.Config.Events {
		if key == name {
			go handler(v)
		}
	}
}
