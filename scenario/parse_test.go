// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/scenario"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailingDefaults(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse", "fail", "bad-defaults.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.ErrorContains(err, "defaults parsing failed")
	assert.Nil(s)
}

func TestNoTests(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "no-tests.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	// When there are plugins but no tests, we should successfully parse the
	// scenario's defaults and have an empty set of Tests in the scenario
	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	assert.Equal("no-tests", s.Name)
	assert.Equal(filepath.Join("testdata", "no-tests.yaml"), s.Path)
	assert.Equal([]string{"books_api", "books_data"}, s.Fixtures)
	assert.Equal(
		map[string]interface{}{
			"foo": &fooDefaults{
				fooInnerDefaults{
					Bar: "barconfig",
				},
			},
			"bar":                &barDefaults{},
			"fail":               &failDefaults{failInnerDefaults{}},
			"priorRun":           &priorRunDefaults{},
			scenario.DefaultsKey: &scenario.Defaults{},
		},
		s.Defaults,
	)
	assert.Empty(s.Tests)
}

func TestFailingPlugin(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "fail-plugin.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.ErrorContains(err, "Indy, bad parse!")
	assert.Nil(s)
}

func TestUnknownSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse", "fail", "unknown-spec.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.NotNil(err)
	assert.Nil(s)
}

func TestBadTimeout(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse", "fail", "bad-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.ErrorIs(err, errors.ErrExpectedMap)
	assert.Nil(s)
}

func TestBadTimeoutDuration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse", "fail", "bad-timeout-duration.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.ErrorContains(err, "invalid duration")
	assert.Nil(s)
}

func TestBadTimeoutDurationScenario(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "parse", "fail", "bad-timeout-duration-scenario.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.ErrorContains(err, "invalid duration")
	assert.Nil(s)
}

func TestKnownSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	assert.Equal("foo", s.Name)
	assert.Equal(filepath.Join("testdata", "foo.yaml"), s.Path)
	assert.Empty(s.Fixtures)
	assert.Equal(
		map[string]interface{}{
			"foo": &fooDefaults{
				fooInnerDefaults{
					Bar: "barconfig",
				},
			},
			"bar":                &barDefaults{},
			"fail":               &failDefaults{failInnerDefaults{}},
			"priorRun":           &priorRunDefaults{},
			scenario.DefaultsKey: &scenario.Defaults{},
		},
		s.Defaults,
	)
	expSpecDefaults := &gdttypes.Defaults{
		"foo": &fooDefaults{
			fooInnerDefaults{
				Bar: "barconfig",
			},
		},
		"bar":                &barDefaults{},
		"fail":               &failDefaults{failInnerDefaults{}},
		"priorRun":           &priorRunDefaults{},
		scenario.DefaultsKey: &scenario.Defaults{},
	}
	expTests := []gdttypes.TestUnit{
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:    0,
				Name:     "bar",
				Defaults: expSpecDefaults,
			},
			Foo: "bar",
		},
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:       1,
				Description: "Bazzy Bizzy",
				Defaults:    expSpecDefaults,
			},
			Foo: "baz",
		},
	}
	assert.Equal(expTests, s.Tests)
}

func TestMultipleSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-bar.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	assert.Equal("foo-bar", s.Name)
	assert.Equal(filepath.Join("testdata", "foo-bar.yaml"), s.Path)
	expTests := []gdttypes.TestUnit{
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:    0,
				Defaults: &gdttypes.Defaults{},
			},
			Foo: "bar",
		},
		&barSpec{
			Spec: gdttypes.Spec{
				Index:    1,
				Defaults: &gdttypes.Defaults{},
			},
			Bar: 42,
		},
	}
	assert.Equal(expTests, s.Tests)
}

func TestEnvExpansion(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "env-expansion.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	t.Setenv("foo", "bar")
	t.Setenv("BAR_CONFIG", "barconfig")
	t.Setenv("DESCRIPTION", "Bazzy Bizzy")

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	assert.Equal("env-expansion", s.Name)
	assert.Equal(filepath.Join("testdata", "env-expansion.yaml"), s.Path)
	assert.Empty(s.Fixtures)
	assert.Equal(
		map[string]interface{}{
			"foo": &fooDefaults{
				fooInnerDefaults{
					Bar: "barconfig",
				},
			},
			"bar":                &barDefaults{},
			"fail":               &failDefaults{failInnerDefaults{}},
			"priorRun":           &priorRunDefaults{},
			scenario.DefaultsKey: &scenario.Defaults{},
		},
		s.Defaults,
	)
	expSpecDefaults := &gdttypes.Defaults{
		"foo": &fooDefaults{
			fooInnerDefaults{
				Bar: "barconfig",
			},
		},
		"bar":                &barDefaults{},
		"fail":               &failDefaults{failInnerDefaults{}},
		"priorRun":           &priorRunDefaults{},
		scenario.DefaultsKey: &scenario.Defaults{},
	}
	expTests := []gdttypes.TestUnit{
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:    0,
				Name:     "$NOT_EXPANDED",
				Defaults: expSpecDefaults,
			},
			Foo: "bar",
		},
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:       1,
				Description: "Bazzy Bizzy",
				Defaults:    expSpecDefaults,
			},
			Foo: "baz",
		},
	}
	assert.Equal(expTests, s.Tests)
}

func TestScenarioDefaults(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	assert.Nil(err)
	assert.NotNil(s)

	assert.IsType(&scenario.Scenario{}, s)
	assert.Equal("foo-timeout", s.Name)
	assert.Equal(filepath.Join("testdata", "foo-timeout.yaml"), s.Path)
	assert.Empty(s.Fixtures)
	assert.Equal(
		map[string]interface{}{
			"foo":      &fooDefaults{},
			"bar":      &barDefaults{},
			"fail":     &failDefaults{failInnerDefaults{}},
			"priorRun": &priorRunDefaults{},
			scenario.DefaultsKey: &scenario.Defaults{
				Timeout: &gdttypes.Timeout{
					After: "2s",
				},
			},
		},
		s.Defaults,
	)
	expSpecDefaults := &gdttypes.Defaults{
		"foo":      &fooDefaults{},
		"bar":      &barDefaults{},
		"fail":     &failDefaults{failInnerDefaults{}},
		"priorRun": &priorRunDefaults{},
		scenario.DefaultsKey: &scenario.Defaults{
			Timeout: &gdttypes.Timeout{
				After: "2s",
			},
		},
	}
	expTests := []gdttypes.TestUnit{
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:    0,
				Defaults: expSpecDefaults,
				Timeout: &gdttypes.Timeout{
					After: "1s",
				},
			},
			Foo: "baz",
		},
		&fooSpec{
			Spec: gdttypes.Spec{
				Index:    1,
				Defaults: expSpecDefaults,
			},
			Foo: "baz",
		},
	}
	assert.Equal(expTests, s.Tests)
}
