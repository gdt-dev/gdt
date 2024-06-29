// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package failer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gdt-dev/gdt/api"
	gdtapi "github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func init() {
	plugin.Register(&Plugin{})
}

type InnerDefaults struct {
	Fail bool `yaml:"fail,omitempty"`
}

type Defaults struct {
	InnerDefaults
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
		case "fail":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			inner := InnerDefaults{}
			if err := valNode.Decode(&inner); err != nil {
				return err
			}
			d.InnerDefaults = inner
			// This is just for testing api when parsing defaults...
			if d.Fail {
				return fmt.Errorf("defaults parsing failed")
			}
		default:
			continue
		}
	}
	return nil
}

type Spec struct {
	api.Spec
	Fail bool `yaml:"fail"`
}

func (s *Spec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *Spec) Base() *api.Spec {
	return &s.Spec
}

func (s *Spec) Retry() *api.Retry {
	return nil
}

func (s *Spec) Timeout() *api.Timeout {
	return nil
}

func (s *Spec) Eval(context.Context) (*api.Result, error) {
	return nil, fmt.Errorf("%w: Indy, bad dates!", gdtapi.RuntimeError)
}

func (s *Spec) UnmarshalYAML(node *yaml.Node) error {
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
		case "fail":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Fail, _ = strconv.ParseBool(valNode.Value)
			if s.Fail {
				return fmt.Errorf("Indy, bad parse!")
			}
		default:
			if lo.Contains(api.BaseSpecFields, key) {
				continue
			}
			return api.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

type Plugin struct{}

func (p *Plugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "fail",
	}
}

func (p *Plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *Plugin) Specs() []api.Evaluable {
	return []api.Evaluable{&Spec{}}
}
