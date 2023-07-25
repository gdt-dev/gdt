// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite_test

import (
	"testing"

	_ "github.com/gdt-dev/gdt/exec"
	"github.com/gdt-dev/gdt/suite"
	"github.com/stretchr/testify/assert"
)

func TestFromDirNoSuchDir(t *testing.T) {
	assert := assert.New(t)

	s, err := suite.FromDir("nosuchdirectory")
	assert.NotNil(err)
	assert.Nil(s)
}

func TestFromDirExecSuite(t *testing.T) {
	assert := assert.New(t)

	s, err := suite.FromDir("testdata/exec")
	assert.Nil(err)
	assert.NotNil(s)

	assert.Equal("testdata/exec", s.Path)
	assert.Len(s.Scenarios, 2)
}
