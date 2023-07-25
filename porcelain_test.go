// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/gdt-dev/gdt/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromUnknownSourceType(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From(1)
	assert.NotNil(err)
	assert.Nil(s)

	assert.ErrorIs(err, gdterrors.ErrUnknownSourceType)
}

func TestFromFileNotFound(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("path", "to", "nonexisting", "file.yaml")
	s, err := gdt.From(fp)
	require.NotNil(err)
	require.Nil(s)

	require.True(os.IsNotExist(err))
}

func TestFromSuite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec")
	s, err := gdt.From(fp)
	require.Nil(err)
	require.NotNil(s)

	suite, ok := s.(*suite.Suite)
	require.True(ok, "gdt.From() did not return a Suite")

	assert.Equal(fp, suite.Path)
	assert.Len(suite.Scenarios, 2)
}

func TestFromScenarioPath(t *testing.T) {
	assert := assert.New(t)

	fp := filepath.Join("suite", "testdata", "exec", "ls.yaml")
	s, err := gdt.From(fp)
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() with dir path did not return a Scenario")

	assert.Equal(fp, sc.Path)
	assert.Len(sc.Tests, 1)
}

func TestFromScenarioReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("suite", "testdata", "exec", "ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)
	s, err := gdt.From(f)
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() from file path did not return a Scenario")

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", sc.Path)
	assert.Len(sc.Tests, 1)
}

func TestFromScenarioBytes(t *testing.T) {
	assert := assert.New(t)

	raw := `name: foo
description: simple foo test
tests:
 - exec: echo foo
`
	b := []byte(raw)
	s, err := gdt.From(b)
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() with []byte did not return a Scenario")

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", sc.Path)
	assert.Len(sc.Tests, 1)
}
