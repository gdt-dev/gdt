// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/plugin"
)

// UnmarshalYAML is a custom unmarshaler that asks plugins for their known spec
// types and attempts to unmarshal test spec contents into those types.
func (s *Scenario) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return api.ExpectedMapAt(node)
	}
	plugins := plugin.Registered()
	defaults := api.Defaults{}
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
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Description = valNode.Value
		case "fixtures":
			if valNode.Kind != yaml.SequenceNode {
				return api.ExpectedSequenceAt(valNode)
			}
			var fixtures []string
			if err := valNode.Decode(&fixtures); err != nil {
				return api.ExpectedSequenceAt(valNode)
			}
			s.Fixtures = fixtures
		case "defaults":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			// Each plugin can have its own set of default configuration values
			// under an outer map field keyed to the name of the plugin.
			// Plugins return a Defaults prototype from
			// `api.Plugin.Defaults()` that understands how to parse a
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
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "tests":
			if valNode.Kind != yaml.SequenceNode {
				return api.ExpectedSequenceAt(valNode)
			}
			for idx, testNode := range valNode.Content {
				parsed := false
				base := api.Spec{}
				if err := testNode.Decode(&base); err != nil {
					return err
				}
				base.Index = idx
				base.Defaults = &defaults
				pluginSpecs := map[api.Plugin][]api.Evaluable{}
				for _, p := range plugins {
					pluginSpecs[p] = p.Specs()
				}
				for plugin, specs := range pluginSpecs {
					for _, sp := range specs {
						if err := testNode.Decode(sp); err != nil {
							if errors.Is(err, api.ErrUnknownField) {
								continue
							}
							return err
						}
						base.Plugin = plugin
						sp.SetBase(base)
						s.Tests = append(s.Tests, sp)
						parsed = true
						break
					}
				}
				if !parsed {
					return api.UnknownSpecAt(s.Path, valNode)
				}
			}
		case "skip-if":
			if valNode.Kind != yaml.SequenceNode {
				return api.ExpectedSequenceAt(valNode)
			}
			for idx, testNode := range valNode.Content {
				parsed := false
				base := api.Spec{}
				if err := testNode.Decode(&base); err != nil {
					return err
				}
				base.Index = idx
				base.Defaults = &defaults
				specs := []api.Evaluable{}
				for _, p := range plugins {
					specs = append(specs, p.Specs()...)
				}
				for _, sp := range specs {
					if err := testNode.Decode(sp); err != nil {
						if errors.Is(err, api.ErrUnknownField) {
							continue
						}
						return err
					}
					sp.SetBase(base)
					s.SkipIf = append(s.SkipIf, sp)
					parsed = true
					break
				}
				if !parsed {
					return api.UnknownSpecAt(s.Path, valNode)
				}
			}
		}
	}
	return nil
}
