// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package fixture

import (
	"context"
	"strings"

	"github.com/gdt-dev/gdt/api"
)

// genericFixture adapts functions and state dicts into the Fixture type
type genericFixture struct {
	starter func(context.Context) error
	stopper func(context.Context)
	state   map[string]interface{}
}

// Start sets up any resources the fixture uses
func (f *genericFixture) Start(ctx context.Context) error {
	if f.starter != nil {
		return f.starter(ctx)
	}
	return nil
}

// Stop cleans up any resources the fixture uses
func (f *genericFixture) Stop(ctx context.Context) {
	if f.stopper != nil {
		f.stopper(ctx)
	}
}

// HasState returns true if the fixture has a state attribute with the supplied
// key
func (f *genericFixture) HasState(key string) bool {
	if f.state != nil {
		_, ok := f.state[strings.ToLower(key)]
		return ok
	}
	return false
}

// State returns a piece of state from the fixture's state map if the supplied
// key exists, otherwise returns nil
func (f *genericFixture) State(key string) interface{} {
	if f.state != nil {
		return f.state[strings.ToLower(key)]
	}
	return nil
}

// genericFixtureModifier sets some value on the test scenario
type genericFixtureModifier func(s *genericFixture)

// WithStarter allows a starter functor to be adapted into a fixture
func WithStarter(starter func(context.Context) error) genericFixtureModifier {
	return func(f *genericFixture) {
		f.starter = starter
	}
}

// WithStopper allows a stopper functor to be adapted into a fixture
func WithStopper(stopper func(context.Context)) genericFixtureModifier {
	return func(f *genericFixture) {
		f.stopper = stopper
	}
}

// WithState allows a map of state key/values to be adapted into a fixture
func WithState(state map[string]interface{}) genericFixtureModifier {
	return func(f *genericFixture) {
		f.state = state
	}
}

// New returns a new generic Fixture
func New(mods ...genericFixtureModifier) api.Fixture {
	f := &genericFixture{}
	for _, mod := range mods {
		mod(f)
	}
	return f
}
