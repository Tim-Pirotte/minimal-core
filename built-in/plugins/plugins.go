package plugins

import (
	"fmt"
	"minimal/minimal-core/built-in/messenger"
	"reflect"
)

type registry struct {
	plugins      plugins
}

type plugins struct {
    messenger *messenger.Messenger
    mapping   map[reflect.Type]Plugin
}

type Plugin interface {
    Init(*plugins)
}

func NewRegistry(messenger *messenger.Messenger) registry {
    return registry{plugins{messenger, map[reflect.Type]Plugin{}}}
}

func (r *registry) Add(plugin Plugin) {
    t := reflect.TypeOf(plugin)

    if t.Kind() != reflect.Pointer {
        r.plugins.logPluginNotAPointer(t)

        return
    }

    if _, ok := r.plugins.mapping[t]; ok {
        r.plugins.logDuplicatePlugin(t)
    }

    r.plugins.mapping[t] = plugin
}

func (r *registry) Setup() {
    for _, plugin := range r.plugins.mapping {
        plugin.Init(&r.plugins)
    }
}

func Get[T any](p *plugins) (T, bool) {
    t := reflect.TypeFor[T]()

    if t.Kind() != reflect.Pointer {
        p.logNotRetrievingPointer(t)
    }

    if p, ok := p.mapping[t]; ok {
        return p.(T), true
    } else {
        var empty T

        return empty, false
    }
}

func (p *plugins) logPluginNotAPointer(t reflect.Type) {
    p.messenger.Send(messenger.Message{
        Message: "Plugin should be a pointer",
        Severity: messenger.Error,
        Notes: []string{
            fmt.Sprintf("Plugin type: %s", t),
            "The plugin wont be added to the registry",
        },
    })
}

func (p *plugins) logDuplicatePlugin(t reflect.Type) {
    p.messenger.Send(messenger.Message{
        Message: "Duplicate plugin declaration",
        Severity: messenger.Error,
        Notes: []string{
            fmt.Sprintf("Plugin type: %s", t),
            "The plugin in the registry will be overwritten",
        },
    })
}

func (p *plugins) logNotRetrievingPointer(t reflect.Type) {
    p.messenger.Send(messenger.Message{
        Message: "Cannot retrieve a non pointer plugin since all plugins are pointers",
        Severity: messenger.Error,
        Notes: []string{fmt.Sprintf("Requested type: %s", t)},
    })
}
