// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt/api"
	gdtexec "github.com/gdt-dev/gdt/plugin/exec"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUnknownShell(t *testing.T) {
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
	assert.ErrorIs(err, api.ErrParse)
	assert.Nil(s)
}

func TestParseSimpleCommand(t *testing.T) {
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
	expTests := []api.Evaluable{
		&gdtexec.Spec{
			Spec: api.Spec{
				Plugin:   gdtexec.PluginRef,
				Index:    0,
				Defaults: &api.Defaults{},
			},
			Action: gdtexec.Action{
				Exec: "ls",
			},
		},
	}
	assert.Equal(expTests, s.Tests)
}

func TestParseVar(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse-var.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	expTests := []api.Evaluable{
		&gdtexec.Spec{
			Spec: api.Spec{
				Plugin:   gdtexec.PluginRef,
				Index:    0,
				Defaults: &api.Defaults{},
			},
			Action: gdtexec.Action{
				Exec: "echo 42",
			},
			Var: gdtexec.Variables{
				"VAR_STDOUT": gdtexec.VarEntry{
					From: "stdout",
				},
				"VAR_STDERR": gdtexec.VarEntry{
					From: "stderr",
				},
				"VAR_RC": gdtexec.VarEntry{
					From: "returncode",
				},
				"MY_ENVVAR": gdtexec.VarEntry{
					From: "MY_ENVVAR",
				},
			},
		},
	}
	assert.Equal(expTests, s.Tests)
}
