// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin_test

import (
	"context"
	"testing"

	"github.com/gdt-dev/gdt/plugin"
	"github.com/gdt-dev/gdt/result"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type fooDefaults struct {
	Foo string `yaml:"foo"`
}

func (d *fooDefaults) UnmarshalYAML(node *yaml.Node) error {
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

func (s *fooSpec) Eval(context.Context) (*result.Result, error) {
	return nil, nil
}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	return nil
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

func TestRegisterAndList(t *testing.T) {
	assert := assert.New(t)

	plugins := plugin.Registered()
	assert.Equal(0, len(plugins))

	plugin.Register(&fooPlugin{})

	plugins = plugin.Registered()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)

	// Add called twice with the same named plugin should be be a no-op

	plugin.Register(&fooPlugin{})

	plugins = plugin.Registered()
	assert.Equal(1, len(plugins))
	assert.Equal("foo", plugins[0].Info().Name)
}
