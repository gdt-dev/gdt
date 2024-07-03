// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
	gdtplugin "github.com/gdt-dev/gdt/plugin"
)

var (
	// this is just for testing purposes...
	PluginRef = &plugin{}
)

func init() {
	gdtplugin.Register(PluginRef)
}

var (
	DefaultTimeout = "10s"
)

// OverrideDefaultTimeout is only used in testing...
func OverrideDefaultTimeout(d string) {
	DefaultTimeout = d
}

const (
	pluginName = "exec"
)

type plugin struct{}

func (p *plugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: pluginName,
		Timeout: &api.Timeout{
			After: DefaultTimeout,
		},
	}
}

func (p *plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *plugin) Specs() []api.Evaluable {
	return []api.Evaluable{&Spec{}}
}

// Plugin returns the HTTP gdt plugin
func Plugin() api.Plugin {
	return &plugin{}
}
