// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context_test

import (
	"context"
	"testing"

	"github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/fixture"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func fooStart(_ context.Context) error { return nil }

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

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

func (s *fooSpec) Eval(ctx context.Context) (*api.Result, error) {
	return nil, nil
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

func TestContext(t *testing.T) {
	assert := assert.New(t)

	ctx := gdtcontext.New()

	assert.Empty(gdtcontext.Plugins(ctx))
	assert.Empty(gdtcontext.Fixtures(ctx))

	ctx = gdtcontext.RegisterPlugin(ctx, &fooPlugin{})
	plugins := gdtcontext.Plugins(ctx)
	assert.Len(plugins, 1)
	assert.Equal("foo", plugins[0].Info().Name)

	fix := fixture.New(fixture.WithStarter(fooStart))
	ctx = gdtcontext.RegisterFixture(ctx, "foo", fix)
	fixtures := gdtcontext.Fixtures(ctx)
	assert.Len(fixtures, 1)
}
