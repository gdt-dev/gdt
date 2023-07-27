// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromUnknownSourceType(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, err := gdt.From(1)
	require.NotNil(err)
	require.Nil(s)

	assert.ErrorIs(err, gdterrors.ErrUnknownSourceType)
}

func TestFromFileNotFound(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("path", "to", "nonexisting", "file.yaml")
	s, err := gdt.From(fp)
	require.NotNil(err)
	require.Nil(s)

	assert.True(os.IsNotExist(err))
}

func TestFromSuite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec")
	suite, err := gdt.From(fp)
	require.Nil(err)
	require.NotNil(suite)

	assert.Equal(fp, suite.Path)
	assert.Len(suite.Scenarios, 2)
}

func TestFromScenarioPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec", "ls.yaml")
	s, err := gdt.From(fp)
	require.Nil(err)
	require.NotNil(s)

	assert.Equal(fp, s.Path)
	assert.Len(s.Scenarios, 1)
	assert.Len(s.Scenarios[0].Tests, 1)
	assert.Equal("exec", s.Name)
}

func TestFromScenarioReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec", "ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)
	suite, err := gdt.From(f)
	assert.Nil(err)
	assert.NotNil(suite)

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", suite.Path)
	assert.Len(suite.Scenarios, 1)
	assert.Len(suite.Scenarios[0].Tests, 1)
}

func TestFromScenarioBytes(t *testing.T) {
	assert := assert.New(t)

	raw := `name: foo
description: simple foo test
tests:
 - exec: echo foo
`
	b := []byte(raw)
	suite, err := gdt.From(b)
	assert.Nil(err)
	assert.NotNil(suite)

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", suite.Path)
	assert.Len(suite.Scenarios, 1)
	assert.Len(suite.Scenarios[0].Tests, 1)
}

func TestRunExecSuite(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec")
	s, err := gdt.From(fp)
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
	require.False(t.Failed())
}

func TestRunExecScenario(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec", "ls.yaml")
	s, err := gdt.From(fp)
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
	require.False(t.Failed())
}
