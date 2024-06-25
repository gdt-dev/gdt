// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
)

type execDefaults struct{}

// Defaults is the known exec plugin defaults collection
type Defaults struct {
	execDefaults
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
		case "exec":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
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
