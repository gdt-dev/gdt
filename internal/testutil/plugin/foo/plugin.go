// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package foo

import (
	"context"
	"fmt"

	"github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/debug"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func init() {
	plugin.Register(&Plugin{})
}

type InnerDefaults struct {
	Bar string `yaml:"bar,omitempty"`
}

type Defaults struct {
	InnerDefaults `yaml:",inline"`
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
		case "foo":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			inner := InnerDefaults{}
			if err := valNode.Decode(&inner); err != nil {
				return err
			}
			d.InnerDefaults = inner
		default:
			continue
		}
	}
	return nil
}

type Spec struct {
	api.Spec
	Foo string `yaml:"foo"`
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
		case "foo":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Foo = valNode.Value
		default:
			if lo.Contains(api.BaseSpecFields, key) {
				continue
			}
			return api.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func (s *Spec) Eval(ctx context.Context) (*api.Result, error) {
	fails := []error{}
	debug.Println(ctx, "in %s Foo=%s", s.Title(), s.Foo)
	// This is just a silly test to demonstrate how to write Eval() methods
	// for plugin Spec specialization classes.
	if s.Name == "bar" && s.Foo != "bar" {
		fail := fmt.Errorf("expected s.Foo = 'bar', got %s", s.Foo)
		fails = append(fails, fail)
	} else if s.Name != "bar" && s.Foo != "baz" {
		fail := fmt.Errorf("expected s.Foo = 'baz', got %s", s.Foo)
		fails = append(fails, fail)
	}
	return api.NewResult(api.WithFailures(fails...)), nil
}

type Plugin struct{}

func (p *Plugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "foo",
	}
}

func (p *Plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *Plugin) Specs() []api.Evaluable {
	return []api.Evaluable{&Spec{}}
}
