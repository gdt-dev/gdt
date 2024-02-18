// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/require"
)

func TestNoExitCodeSimpleCommand(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestExitCode(t *testing.T) {
	require := require.New(t)

	fname := "ls-with-exit-code.yaml"
	// Yay, different exit codes for the same not found error...
	if runtime.GOOS == "darwin" {
		fname = "mac-ls-with-exit-code.yaml"
	}

	fp := filepath.Join("testdata", fname)
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestShellList(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "shell-ls.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestIs(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "echo-cat.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestContains(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-contains.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestContainsOneOf(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-contains-one-of.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestContainsNoneOf(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-contains-none-of.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestSleepTimeout(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "sleep-timeout.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := context.TODO()
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestDebugWriter(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "echo-cat.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))
	err = s.Run(ctx, t)
	require.Nil(err)
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "exec: echo [cat]")
	require.Contains(debugout, "exec: stdout: cat")
	require.Contains(debugout, "exec: sh [-c echo cat 1>&2]")
	require.Contains(debugout, "exec: stderr: cat")
}

func TestWait(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "echo-wait.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))
	err = s.Run(ctx, t)
	require.Nil(err)
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "wait: 10ms before")
	require.Contains(debugout, "wait: 20ms after")
}

func TestTimeoutCascade(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "timeout-cascade.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	ctx := gdtcontext.New(gdtcontext.WithDebug(w))
	err = s.Run(ctx, t)
	require.Nil(err)
	require.False(t.Failed())
	w.Flush()
	require.NotEqual(b.Len(), 0)
	debugout := b.String()
	require.Contains(debugout, "using timeout of 500ms (expected: false) [scenario default]")
	require.Contains(debugout, "using timeout of 20ms (expected: true)")
}

// Unfortunately there's not really any good way of testing things like this
// except by manually causing an assertion to fail in the test case and
// checking to see if the `on.fail` action was taken and debug output emitted
// to the console.
//
// When I change the `testdata/on-fail-exec.yaml` file to have a failed
// assertion by changing `assert.out.is` to "dat" instead of "cat", I get the
// correct behaviour:
//
// === RUN   TestOnFail
// === RUN   TestOnFail/on-fail-exec
//
//	action.go:59: exec: echo [cat]
//	eval.go:35: assertion failed: not equal: expected dat but got cat
//	action.go:59: exec: echo [bad kitty]
//	eval.go:46: on.fail.exec: stdout: bad kitty
//
// === NAME  TestOnFail
//
//	eval_test.go:256:
//	    	Error Trace:	/home/jaypipes/src/github.com/gdt-dev/gdt/plugin/exec/eval_test.go:256
//	    	Error:      	Should be false
//	    	Test:       	TestOnFail
//
// --- FAIL: TestOnFail (0.00s)
//
//	--- FAIL: TestOnFail/on-fail-exec (0.00s)
func TestOnFail(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "on-fail-exec.yaml")
	f, err := os.Open(fp)
	require.Nil(err)

	s, err := scenario.FromReader(
		f,
		scenario.WithPath(fp),
	)
	require.Nil(err)
	require.NotNil(s)

	ctx := gdtcontext.New(gdtcontext.WithDebug())
	err = s.Run(ctx, t)
	require.Nil(err)
	require.False(t.Failed())
}
