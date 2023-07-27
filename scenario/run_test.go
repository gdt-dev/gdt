// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/result"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *fooSpec) Eval(ctx context.Context, t *testing.T) *result.Result {
	fails := []error{}
	t.Run(s.Title(), func(t *testing.T) {
		debug.Printf(ctx, t, "in %s Foo=%s", s.Title(), s.Foo)
		// This is just a silly test to demonstrate how to write Eval() methods
		// for plugin Spec specialization classes.
		if s.Name == "bar" && s.Foo != "bar" {
			fail := fmt.Errorf("expected s.Foo = 'bar', got %s", s.Foo)
			fails = append(fails, fail)
		} else if s.Name != "bar" && s.Foo != "baz" {
			fail := fmt.Errorf("expected s.Foo = 'baz', got %s", s.Foo)
			fails = append(fails, fail)
		}
	})
	return result.New(result.WithFailures(fails...))
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
	assert.ErrorIs(err, gdterrors.ErrRequiredFixture)
	assert.ErrorIs(err, gdterrors.RuntimeError)
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
