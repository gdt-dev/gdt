// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package fixture_test

import (
	"testing"

	"github.com/gdt-dev/gdt/fixture"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	assert := assert.New(t)

	state := map[string]interface{}{
		"foo": "bar",
		"baz": 1,
	}

	f := fixture.New(
		fixture.WithState(state),
	)

	assert.True(f.HasState("foo"))
	assert.Equal("bar", f.State("foo"))
	assert.True(f.HasState("baz"))
	assert.Equal(1, f.State("baz"))
	assert.False(f.HasState("bar"))
}

func TestStarter(t *testing.T) {
	assert := assert.New(t)

	started := false

	starter := func() {
		started = true
	}

	f := fixture.New(
		fixture.WithStarter(starter),
	)

	assert.False(started)

	f.Start()

	assert.True(started)
}

func TestStopper(t *testing.T) {
	assert := assert.New(t)

	stopped := false

	stopper := func() {
		stopped = true
	}

	f := fixture.New(
		fixture.WithStopper(stopper),
	)

	assert.False(stopped)

	f.Stop()

	assert.True(stopped)
}
