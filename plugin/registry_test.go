// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package plugin_test

import (
	"context"
	"testing"

	"github.com/gdt-dev/gdt/api"
	"github.com/gdt-dev/gdt/plugin"
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

func (s *fooSpec) Eval(context.Context) (*api.Result, error) {
	return nil, nil
}

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	return nil
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
