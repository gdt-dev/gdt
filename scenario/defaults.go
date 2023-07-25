// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"time"

	"github.com/gdt-dev/gdt/errors"
	gdttypes "github.com/gdt-dev/gdt/types"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultsKey is the key within the Defaults collection for
	// scenario defaults.
	DefaultsKey = "gdt.scenario"
)

// Defaults is the scenario's defaults collection
type Defaults struct {
	// Timeout has fields that represent the default timeout behaviour and
	// expectations to use for test specs in the scenario.
	Timeout *gdttypes.Timeout `yaml:"timeout,omitempty"`
}

func (d *Defaults) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "timeout":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			var to *gdttypes.Timeout
			if err := valNode.Decode(&to); err != nil {
				return errors.ExpectedTimeoutAt(valNode)
			}
			_, err := time.ParseDuration(to.After)
			if err != nil {
				return err
			}
			d.Timeout = to
		default:
			continue
		}
	}
	return nil
}
