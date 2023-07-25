// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite_test

import (
	"testing"

	"github.com/gdt-dev/gdt/suite"
	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	assert := assert.New(t)

	s := suite.New(
		suite.WithPath("/path/to/suite"),
	)

	assert.Equal("/path/to/suite", s.Path)
}
