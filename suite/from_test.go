// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite_test

import (
	"testing"

	_ "github.com/gdt-dev/gdt/plugin/exec"
	"github.com/gdt-dev/gdt/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromDirNoSuchDir(t *testing.T) {
	require := require.New(t)

	s, err := suite.FromDir("nosuchdirectory")
	require.NotNil(err)
	require.Nil(s)
}

func TestFromDirExecSuite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, err := suite.FromDir("testdata/exec")
	require.Nil(err)
	require.NotNil(s)

	assert.Equal("testdata/exec", s.Path)
	assert.Len(s.Scenarios, 2)
}
