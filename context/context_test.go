// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package context_test

import (
	"context"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/fixture"
	"github.com/gdt-dev/gdt/result"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func fooStart(_ context.Context) {}

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

func (s *fooSpec) UnmarshalYAML(node *yaml.Node) error {
	return nil
}

func (s *fooSpec) Eval(ctx context.Context, t *testing.T) *result.Result {
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
