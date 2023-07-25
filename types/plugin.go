// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import "gopkg.in/yaml.v3"

// PluginInfo contains basic information about the plugin and what type of
// tests it can handle.
type PluginInfo struct {
	// Name is the primary name of the plugin
	Name string
	// Aliases is an optional set of aliased names for the plugin
	Aliases []string
	// Description describes what types of tests the plugin can handle.
	Description string
}

// Plugin is the driver interface for different types of gdt tests.
type Plugin interface {
	// Info returns a struct that describes what the plugin does
	Info() PluginInfo
	// Defaults returns a YAML Unmarshaler types that the plugin knows how
	// to parse its defaults configuration with.
	Defaults() yaml.Unmarshaler
	// Specs returns a list of YAML Unmarshaler types that the plugin knows
	// how to parse.
	Specs() []TestUnit
}
