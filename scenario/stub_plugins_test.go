// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	"github.com/gdt-dev/gdt/errors"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/plugin"
	"github.com/gdt-dev/gdt/result"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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
		case "fail":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			inner := failInnerDefaults{}
			if err := valNode.Decode(&inner); err != nil {
				return err
			}
			d.failInnerDefaults = inner
			// This is just for testing errors when parsing defaults...
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
	gdttypes.Spec
	Fail bool `yaml:"fail"`
}

func (s *failSpec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *failSpec) Base() *gdttypes.Spec {
	return &s.Spec
}

func (s *failSpec) Eval(context.Context, *testing.T) *result.Result {
	return result.New(
		result.WithRuntimeError(
			fmt.Errorf("%w: Indy, bad dates!", gdterrors.RuntimeError),
		),
	)
}

func (s *failSpec) UnmarshalYAML(node *yaml.Node) error {
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
		case "fail":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Fail, _ = strconv.ParseBool(valNode.Value)
			if s.Fail {
				return fmt.Errorf("Indy, bad parse!")
			}
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

type failingPlugin struct{}

func (p *failingPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "fail",
	}
}

func (p *failingPlugin) Defaults() yaml.Unmarshaler {
	return &failDefaults{}
}

func (p *failingPlugin) Specs() []gdttypes.Evaluable {
	return []gdttypes.Evaluable{&failSpec{}}
}

type fooInnerDefaults struct {
	Bar string `yaml:"bar,omitempty"`
}

type fooDefaults struct {
	fooInnerDefaults
}

func (d *fooDefaults) UnmarshalYAML(node *yaml.Node) error {
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
		case "foo":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
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
	gdttypes.Spec
	Foo string `yaml:"foo"`
}

func (s *fooSpec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *fooSpec) Base() *gdttypes.Spec {
	return &s.Spec
}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
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
		case "foo":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Foo = valNode.Value
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func (s *fooSpec) Eval(ctx context.Context, t *testing.T) *result.Result {
	fails := []error{}
	t.Run(s.Title(), func(t *testing.T) {
		debug.Printf(ctx, t, "in %s Foo=%s", s.Title(), s.Foo)
		// This is just a silly test to demonstrate how to write Eval() methods
		// for plugin Spec specialization classes.
		if s.Name == "bar" && s.Foo != "bar" {
			fail := fmt.Errorf("expected s.Foo = 'bar', got %s", s.Foo)
			fails = append(fails, fail)
		} else if s.Name != "bar" && s.Foo != "baz" {
			fail := fmt.Errorf("expected s.Foo = 'baz', got %s", s.Foo)
			fails = append(fails, fail)
		}
	})
	return result.New(result.WithFailures(fails...))
}

type fooPlugin struct{}

func (p *fooPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "foo",
	}
}

func (p *fooPlugin) Defaults() yaml.Unmarshaler {
	return &fooDefaults{}
}

func (p *fooPlugin) Specs() []gdttypes.Evaluable {
	return []gdttypes.Evaluable{&fooSpec{}}
}

type barDefaults struct {
	Foo string `yaml:"bar"`
}

func (d *barDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type barSpec struct {
	gdttypes.Spec
	Bar int `yaml:"bar"`
}

func (s *barSpec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *barSpec) Base() *gdttypes.Spec {
	return &s.Spec
}

func (s *barSpec) Eval(context.Context, *testing.T) *result.Result {
	return result.New()
}

func (s *barSpec) UnmarshalYAML(node *yaml.Node) error {
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
		case "bar":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			if v, err := strconv.Atoi(valNode.Value); err != nil {
				return errors.ExpectedIntAt(valNode)
			} else {
				s.Bar = v
			}
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

type barPlugin struct{}

func (p *barPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "bar",
	}
}

func (p *barPlugin) Defaults() yaml.Unmarshaler {
	return &barDefaults{}
}

func (p *barPlugin) Specs() []gdttypes.Evaluable {
	return []gdttypes.Evaluable{&barSpec{}}
}

const priorRunDataKey = "priorrun"

type priorRunDefaults struct{}

func (d *priorRunDefaults) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

type priorRunSpec struct {
	gdttypes.Spec
	State string `yaml:"state"`
	Prior string `yaml:"prior"`
}

func (s *priorRunSpec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *priorRunSpec) Base() *gdttypes.Spec {
	return &s.Spec
}

func (s *priorRunSpec) UnmarshalYAML(node *yaml.Node) error {
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
		case "state":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.State = valNode.Value
		case "prior":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Prior = valNode.Value
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func (s *priorRunSpec) Eval(ctx context.Context, t *testing.T) *result.Result {
	t.Run(s.Title(), func(t *testing.T) {
		assert := assert.New(t)
		// Here we test that the prior run data that we save at the end of each
		// Run() is showing up properly in the next Run()'s context.
		prData := gdtcontext.PriorRun(ctx)
		if s.Index == 0 {
			assert.Empty(prData)
		} else {
			assert.Contains(prData, priorRunDataKey)
			assert.IsType(prData[priorRunDataKey], "")
			assert.Equal(s.Prior, prData[priorRunDataKey])
		}
	})
	return result.New(result.WithData(priorRunDataKey, s.State))
}

type priorRunPlugin struct{}

func (p *priorRunPlugin) Info() gdttypes.PluginInfo {
	return gdttypes.PluginInfo{
		Name: "priorRun",
	}
}

func (p *priorRunPlugin) Defaults() yaml.Unmarshaler {
	return &priorRunDefaults{}
}

func (p *priorRunPlugin) Specs() []gdttypes.Evaluable {
	return []gdttypes.Evaluable{&priorRunSpec{}}
}
