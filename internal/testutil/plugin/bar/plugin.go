// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package bar

import (
	"context"
	"strconv"

	"github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

var (
	// this is just for testing purposes...
	PluginRef = &Plugin{}
)

func init() {
	plugin.Register(PluginRef)
}

type Defaults struct {
	Foo string `yaml:"bar"`
}

func (d *Defaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type Spec struct {
	api.Spec
	Bar int `yaml:"bar"`
}

func (s *Spec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *Spec) Base() *api.Spec {
	return &s.Spec
}

func (s *Spec) Retry() *api.Retry {
	return api.NoRetry
}

func (s *Spec) Timeout() *api.Timeout {
	return nil
}

func (s *Spec) Eval(context.Context) (*api.Result, error) {
	return api.NewResult(), nil
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
		case "bar":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			if v, err := strconv.Atoi(valNode.Value); err != nil {
				return api.ExpectedIntAt(valNode)
			} else {
				s.Bar = v
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
		Name: "bar",
	}
}

func (p *Plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *Plugin) Specs() []api.Evaluable {
	return []api.Evaluable{&Spec{}}
}
