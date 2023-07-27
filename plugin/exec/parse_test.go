// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt/errors"
	gdtexec "github.com/gdt-dev/gdt/plugin/exec"
	"github.com/gdt-dev/gdt/scenario"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnknownShell(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "unknown-shell-test.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	assert.NotNil(err)
	assert.ErrorIs(err, gdtexec.ErrUnknownShell)
	assert.ErrorIs(err, errors.ErrParse)
	assert.Nil(s)
}

func TestSimpleCommand(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	expTests := []gdttypes.Evaluable{
		&gdtexec.Spec{
			Spec: gdttypes.Spec{
				Index:    0,
				Defaults: &gdttypes.Defaults{},
			},
			Exec: "ls",
		},
	}
	assert.Equal(expTests, s.Tests)
}
