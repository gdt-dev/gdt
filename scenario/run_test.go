// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario_test

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gdt-dev/gdt/internal/testutil/fixture/errstarter"
)

var failFlag = flag.Bool("fail", false, "run tests expected to fail")

func TestRun(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "foo.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
}

func TestPriorRun(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "prior-run.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
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
	assert.ErrorIs(err, api.ErrRequiredFixture)
	assert.ErrorIs(err, api.RuntimeError)
}

func TestTimeoutConflictTotalWait(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	fp := filepath.Join("testdata", "timeout-conflict-total-wait.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	assert.NotNil(err)
	assert.ErrorIs(err, api.ErrTimeoutConflict)
	assert.ErrorIs(err, api.RuntimeError)
}

func TestTimeoutConflictSpecTimeout(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	fp := filepath.Join("testdata", "timeout-conflict-spec-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	assert.NotNil(err)
	assert.ErrorIs(err, api.ErrTimeoutConflict)
	assert.ErrorIs(err, api.RuntimeError)
}

func TestTimeoutConflictDefaultTimeout(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	fp := filepath.Join("testdata", "timeout-conflict-default-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	assert.NotNil(err)
	assert.ErrorIs(err, api.ErrTimeoutConflict)
	assert.ErrorIs(err, api.RuntimeError)
}

func TestFixtureStartError(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	fp := filepath.Join("testdata", "fixture-start-error.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	ctx := gdtcontext.New()
	ctx = gdtcontext.RegisterFixture(ctx, "start-error", errstarter.Fixture)

	err = s.Run(ctx, t)
	assert.NotNil(err)
	assert.ErrorContains(err, "error starting fixture!")
}

func TestDebugFlushing(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "foo-debug-wait-flush.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(ctx, t)
	require.Nil(err)
	require.False(t.Failed())
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "[gdt] [foo-debug-wait-flush/0:bar] wait: 250ms before")
}

func TestNoRetry(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "no-retry.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(ctx, t)
	require.Nil(err)
	require.False(t.Failed())
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "[gdt] [no-retry/0:bar] spec/run: single-shot (no retries) ok: true")
}

func TestNoRetryEvaluableOverride(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "no-retry-evaluable-override.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(ctx, t)
	require.Nil(err)
	require.False(t.Failed())
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "[gdt] [no-retry-evaluable-override/0:bar] spec/run: single-shot (no retries) ok: true")
}

func TestFailRetryTestOverride(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
	require := require.New(t)

	fp := filepath.Join("testdata", "retry-test-override.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	ctx := gdtcontext.New(gdtcontext.WithDebug())
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestRetryTestOverride(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailRetryTestOverride",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)

	debugout := string(outerr)
	require.Contains(debugout, "[gdt] [retry-test-override/0:baz] spec/run: exceeded max attempts 2. stopping.")
}

func TestSkipIf(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "skip-if.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(f, scenario.WithPath(fp))
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(context.TODO(), t)
	require.Nil(err)
	require.True(t.Skipped())
}
