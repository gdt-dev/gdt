// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"testing"

	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	assert := assert.New(t)

	s := scenario.New(
		scenario.WithPath("/path/to/foo.yaml"),
	)

	assert.Equal("/path/to/foo.yaml", s.Path)
	assert.Equal("", s.Name)
	// Title() returns the basename from the path if name isn't present
	assert.Equal("foo.yaml", s.Title())

	s.Name = "foo"
	assert.Equal("foo", s.Title())
}
