// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin

import (
	"strings"
	"sync"

	"github.com/gdt-dev/gdt/api"
)

// registry stores a set of Plugins and is safe to use in threaded
// environments.
type registry struct {
	sync.RWMutex
	entries map[string]api.Plugin
}

// Remove delists the Plugin with registry. Only really useful for testing.
func (r *registry) Remove(p api.Plugin) {
	r.Lock()
	defer r.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	delete(r.entries, lowered)
}

// Add registers a Plugin with the registry.
func (r *registry) Add(p api.Plugin) {
	r.Lock()
	defer r.Unlock()
	lowered := strings.ToLower(p.Info().Name)
	r.entries[lowered] = p
}

// List returns a slice of Plugins that are registered with gdt.
func (r *registry) List() []api.Plugin {
	r.RLock()
	defer r.RUnlock()
	res := []api.Plugin{}
	for _, p := range r.entries {
		res = append(res, p)
	}
	return res
}

var (
	knownPlugins = &registry{
		entries: map[string]api.Plugin{},
	}
)

// Register registers a plugin with gdt's set of known plugins.
//
// Generally only plugin authors will ever need to call this function. It is
// not required for normal use of gdt or any known plugin.
func Register(p api.Plugin) {
	knownPlugins.Add(p)
}

// Registered returns a slice of pointers to gdt's known plugins.
func Registered() []api.Plugin {
	return knownPlugins.List()
}
