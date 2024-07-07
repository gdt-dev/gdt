// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"time"

	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
)

const (
	// DefaultsKey is the key within the Defaults collection for
	// scenario defaults. Note that this isn't exposed in the YAML schema for a
	// scenario. It's just used as a way of indicating to the scenario runner
	// what was found in the scenario YAML's `defaults` top-level field.
	DefaultsKey = "gdt.scenario"
)

// Defaults is the scenario's defaults collection
type Defaults struct {
	// Timeout has fields that represent the default timeout behaviour and
	// expectations to use for test specs in the scenario.
	Timeout *api.Timeout `yaml:"timeout,omitempty"`
	// Retry has fields that represent the default retry behaviour for test
	// specs in the scenario.
	Retry *api.Retry `yaml:"retry,omitempty"`
}

func (d *Defaults) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return api.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "timeout":
			var to *api.Timeout
			switch valNode.Kind {
			case yaml.MappingNode:
				// We support the old-style timeout:after
				if err := valNode.Decode(&to); err != nil {
					return api.ExpectedTimeoutAt(valNode)
				}
			case yaml.ScalarNode:
				// We also support a straight string duration
				to = &api.Timeout{
					After: valNode.Value,
				}
			default:
				return api.ExpectedScalarOrMapAt(valNode)
			}
			_, err := time.ParseDuration(to.After)
			if err != nil {
				return err
			}
			d.Timeout = to
		case "retry":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			var r *api.Retry
			if err := valNode.Decode(&r); err != nil {
				return api.ExpectedRetryAt(valNode)
			}
			if r.Attempts != nil {
				attempts := *r.Attempts
				if attempts < 1 {
					return api.InvalidRetryAttempts(valNode, attempts)
				}
			}
			if r.Interval != "" {
				_, err := time.ParseDuration(r.Interval)
				if err != nil {
					return err
				}
			}
			d.Retry = r
		default:
			continue
		}
	}
	return nil
}
