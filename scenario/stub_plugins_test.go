// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gdt-dev/gdt/api"
	gdtapi "github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

func init() {
	plugin.Register(&fooPlugin{})
	plugin.Register(&barPlugin{})
	plugin.Register(&failingPlugin{})
	plugin.Register(&priorRunPlugin{})
}

type failInnerDefaults struct {
	Fail bool `yaml:"fail,omitempty"`
}

type failDefaults struct {
	failInnerDefaults
}

func (d *failDefaults) UnmarshalYAML(node *yaml.Node) error {
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
			inner := failInnerDefaults{}
			if err := valNode.Decode(&inner); err != nil {
				return err
			}
			d.failInnerDefaults = inner
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

type failSpec struct {
	api.Spec
	Fail bool `yaml:"fail"`
}

func (s *failSpec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *failSpec) Base() *api.Spec {
	return &s.Spec
}

func (s *failSpec) Retry() *api.Retry {
	return nil
}

func (s *failSpec) Timeout() *api.Timeout {
	return nil
}

func (s *failSpec) Eval(context.Context) (*api.Result, error) {
	return nil, fmt.Errorf("%w: Indy, bad dates!", gdtapi.RuntimeError)
}

func (s *failSpec) UnmarshalYAML(node *yaml.Node) error {
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

type failingPlugin struct{}

func (p *failingPlugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "fail",
	}
}

func (p *failingPlugin) Defaults() yaml.Unmarshaler {
	return &failDefaults{}
}

func (p *failingPlugin) Specs() []api.Evaluable {
	return []api.Evaluable{&failSpec{}}
}

type fooInnerDefaults struct {
	Bar string `yaml:"bar,omitempty"`
}

type fooDefaults struct {
	fooInnerDefaults
}

func (d *fooDefaults) UnmarshalYAML(node *yaml.Node) error {
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
			inner := fooInnerDefaults{}
			if err := valNode.Decode(&inner); err != nil {
				return err
			}
			d.fooInnerDefaults = inner
		default:
			continue
		}
	}
	return nil
}

type fooSpec struct {
	api.Spec
	Foo string `yaml:"foo"`
}

func (s *fooSpec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *fooSpec) Base() *api.Spec {
	return &s.Spec
}

func (s *fooSpec) Retry() *api.Retry {
	return nil
}

func (s *fooSpec) Timeout() *api.Timeout {
	return nil
}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
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

func (s *fooSpec) Eval(ctx context.Context) (*api.Result, error) {
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

type fooPlugin struct{}

func (p *fooPlugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "foo",
	}
}

func (p *fooPlugin) Defaults() yaml.Unmarshaler {
	return &fooDefaults{}
}

func (p *fooPlugin) Specs() []api.Evaluable {
	return []api.Evaluable{&fooSpec{}}
}

type barDefaults struct {
	Foo string `yaml:"bar"`
}

func (d *barDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type barSpec struct {
	api.Spec
	Bar int `yaml:"bar"`
}

func (s *barSpec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *barSpec) Base() *api.Spec {
	return &s.Spec
}

func (s *barSpec) Retry() *api.Retry {
	return api.NoRetry
}

func (s *barSpec) Timeout() *api.Timeout {
	return nil
}

func (s *barSpec) Eval(context.Context) (*api.Result, error) {
	return api.NewResult(), nil
}

func (s *barSpec) UnmarshalYAML(node *yaml.Node) error {
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

type barPlugin struct{}

func (p *barPlugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "bar",
	}
}

func (p *barPlugin) Defaults() yaml.Unmarshaler {
	return &barDefaults{}
}

func (p *barPlugin) Specs() []api.Evaluable {
	return []api.Evaluable{&barSpec{}}
}

const priorRunDataKey = "priorrun"

type priorRunDefaults struct{}

func (d *priorRunDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type priorRunSpec struct {
	api.Spec
	State string `yaml:"state"`
	Prior string `yaml:"prior"`
}

func (s *priorRunSpec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *priorRunSpec) Base() *api.Spec {
	return &s.Spec
}

func (s *priorRunSpec) Retry() *api.Retry {
	return nil
}

func (s *priorRunSpec) Timeout() *api.Timeout {
	return nil
}

func (s *priorRunSpec) UnmarshalYAML(node *yaml.Node) error {
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

func (s *priorRunSpec) Eval(ctx context.Context) (*api.Result, error) {
	// Here we test that the prior run data that we save at the end of each
	// Run() is showing up properly in the next Run()'s context.
	fails := []error{}
	prData := gdtcontext.PriorRun(ctx)
	if s.Index == 0 {
		if len(prData) != 0 {
			fails = append(fails, fmt.Errorf("expected prData to be empty"))
		}
	} else {
		data, ok := prData[priorRunDataKey]
		if !ok {
			fails = append(fails, fmt.Errorf("expected priorRunDataKey in priorRun map"))
		}
		if s.Prior != data {
			fails = append(fails, fmt.Errorf("expected priorRunData == data"))
		}
	}
	return api.NewResult(
		api.WithFailures(fails...),
		api.WithData(priorRunDataKey, s.State),
	), nil
}

type priorRunPlugin struct{}

func (p *priorRunPlugin) Info() api.PluginInfo {
	return api.PluginInfo{
		Name: "priorRun",
	}
}

func (p *priorRunPlugin) Defaults() yaml.Unmarshaler {
	return &priorRunDefaults{}
}

func (p *priorRunPlugin) Specs() []api.Evaluable {
	return []api.Evaluable{&priorRunSpec{}}
}
