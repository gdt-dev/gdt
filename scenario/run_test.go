// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *fooSpec) Run(ctx context.Context, t *testing.T) error {
	t.Run(s.Title(), func(t *testing.T) {
		assert := assert.New(t)
		debug.Printf(ctx, t, "in %s Foo=%s", s.Title(), s.Foo)
		// This is just a silly test to demonstrate how to write Run() commands
		// for plugin Spec specialization classes.
		if s.Name == "bar" {
			assert.Equal(s.Foo, "bar")
		} else {
			assert.Equal(s.Foo, "baz")
		}
	})
	return nil
}

func TestRun(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "foo.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	s.Run(context.TODO(), t)
}

func TestPriorRun(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "prior-run.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	s.Run(context.TODO(), t)
}

func TestMissingFixtures(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	fp := filepath.Join("testdata", "foo-fixtures.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	// Pass a context with no fixtures registered...
	err = s.Run(context.TODO(), t)
	assert.NotNil(err)
	assert.ErrorIs(err, gdterrors.ErrRuntime)
	assert.ErrorIs(err, gdterrors.ErrRequiredFixture)
}

func TestDebugFlushing(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-debug-wait-flush.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	ctx := gdtcontext.New(
		gdtcontext.WithDebug(),
	)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestTimeoutCascade(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
}
