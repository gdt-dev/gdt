// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/errors"
	gdttypes "github.com/gdt-dev/gdt/types"
)

type execDefaults struct{}

// Defaults is the known exec plugin defaults collection
type Defaults struct {
	execDefaults
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
		case "exec":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			ed := execDefaults{}
			if err := valNode.Decode(&ed); err != nil {
				return err
			}
			d.execDefaults = ed
		default:
			continue
		}
	}
	return nil
}

// fromBaseDefaults returns an exec plugin-specific Defaults from a Spec
func fromBaseDefaults(base *gdttypes.Defaults) *Defaults {
	if base == nil {
		return nil
	}
	d := base.For(pluginName)
	if d == nil {
		return nil
	}
	return d.(*Defaults)
}
