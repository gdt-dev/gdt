// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package priorrun

import (
	"context"
	"fmt"

	"github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func init() {
	plugin.Register(&Plugin{})
}

const PriorRunDataKey = "priorrun"

type Defaults struct{}

func (d *Defaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type Spec struct {
	api.Spec
	State string `yaml:"state"`
	Prior string `yaml:"prior"`
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
		case "state":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.State = valNode.Value
		case "prior":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Prior = valNode.Value
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
	// Here we test that the prior run data that we save at the end of each
	// Run() is showing up properly in the next Run()'s context.
	fails := []error{}
	prData := gdtcontext.PriorRun(ctx)
	if s.Index == 0 {
		if len(prData) != 0 {
			fails = append(fails, fmt.Errorf("expected prData to be empty"))
		}
	} else {
		data, ok := prData[PriorRunDataKey]
		if !ok {
			fails = append(fails, fmt.Errorf("expected PriorRunDataKey in priorRun map"))
		}
		if s.Prior != data {
			fails = append(fails, fmt.Errorf("expected priorRunData == data"))
		}
	}
	return api.NewResult(
		api.WithFailures(fails...),
		api.WithData(PriorRunDataKey, s.State),
	), nil
}

type Plugin struct{}

func (p *Plugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "priorRun",
	}
}

func (p *Plugin) Defaults() yaml.Unmarshaler {
	return &Defaults{}
}

func (p *Plugin) Specs() []api.Evaluable {
	return []api.Evaluable{&Spec{}}
}
