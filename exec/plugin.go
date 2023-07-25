// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"gopkg.in/yaml.v3"

	gdtplugin "github.com/gdt-dev/gdt/plugin"
	gdttypes "github.com/gdt-dev/gdt/types"
)

func init() {
	gdtplugin.Register(Plugin())
}

const (
	pluginName = "exec"
)

type plugin struct{}

func (p *plugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: pluginName,
	}
}

func (p *plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *plugin) Specs() []gdttypes.TestUnit {
	return []gdttypes.TestUnit{&Spec{}}
}

// Plugin returns the HTTP gdt plugin
func Plugin() gdttypes.Plugin {
	return &plugin{}
}
