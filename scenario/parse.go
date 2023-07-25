// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"errors"

	"gopkg.in/yaml.v3"

	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/plugin"
	gdttypes "github.com/gdt-dev/gdt/types"
)

// UnmarshalYAML is a custom unmarshaler that asks plugins for their known spec
// types and attempts to unmarshal test spec contents into those types.
func (s *Scenario) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return gdterrors.ExpectedMapAt(node)
	}
	plugins := plugin.Registered()
	defaults := gdttypes.Defaults{}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	//
	// We do a first pass over the scenario's common fields in order to get a
	// set of parsed defaults and required fixtures. We then parse the
	// individual test units since those are plugin-specific and may rely on
	// the parsed defaults and required fixtures.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return gdterrors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return gdterrors.ExpectedScalarAt(valNode)
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return gdterrors.ExpectedScalarAt(valNode)
			}
			s.Description = valNode.Value
		case "require":
			if valNode.Kind != yaml.SequenceNode {
				return gdterrors.ExpectedSequenceAt(valNode)
			}
			requires := make([]string, len(valNode.Content))
			for x, n := range valNode.Content {
				requires[x] = n.Value
			}
			s.Require = requires
		case "defaults":
			if valNode.Kind != yaml.MappingNode {
				return gdterrors.ExpectedMapAt(valNode)
			}
			// Each plugin can have its own set of default configuration values
			// under an outer map field keyed to the name of the plugin.
			// Plugins return a Defaults prototype from
			// `gdttypes.Plugin.Defaults()` that understands how to parse a
			// `yaml.Node` that represents the top-level defaults object in the
			// scenario.
			for _, p := range plugins {
				plugDefaults := p.Defaults()
				if err := valNode.Decode(plugDefaults); err != nil {
					return err
				}
				defaults[p.Info().Name] = plugDefaults
			}
			// The scenario may have its own defaults as well, so we stash
			// these in the "scenario" pseudo-plugin key.
			var scenDefaults Defaults
			if err := valNode.Decode(&scenDefaults); err != nil {
				return err
			}
			defaults[DefaultsKey] = &scenDefaults
			s.Defaults = defaults
		}
	}
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return gdterrors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		if key == "tests" {
			valNode := node.Content[i+1]
			if valNode.Kind != yaml.SequenceNode {
				return gdterrors.ExpectedSequenceAt(valNode)
			}
			for idx, testNode := range valNode.Content {
				parsed := false
				base := gdttypes.Spec{}
				if err := testNode.Decode(&base); err != nil {
					return err
				}
				base.Index = idx
				base.Defaults = &defaults
				specs := []gdttypes.TestUnit{}
				for _, p := range plugins {
					specs = append(specs, p.Specs()...)
				}
				for _, sp := range specs {
					if err := testNode.Decode(sp); err != nil {
						if errors.Is(err, gdterrors.ErrUnknownField) {
							continue
						}
						return err
					}
					sp.SetBase(base)
					s.Tests = append(s.Tests, sp)
					parsed = true
					break
				}
				if !parsed {
					return gdterrors.UnknownSpecAt(s.Path, valNode)
				}
			}
		}
	}
	return nil
}
