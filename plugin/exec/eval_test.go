// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec_test

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	gdtcontext "github.com/gdt-dev/gdt/context"
	execplugin "github.com/gdt-dev/gdt/plugin/exec"
	"github.com/gdt-dev/gdt/scenario"
	"github.com/stretchr/testify/require"
)

func init() {
	execplugin.OverrideDefaultTimeout("0.5s")
}

var failFlag = flag.Bool("fail", false, "run tests expected to fail")

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

func TestFailExecExitCodeNotSpecified(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
	require := require.New(t)

	fp := filepath.Join("testdata", "ls-fail-no-exit-code.yaml")
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
}

func TestExecFailExitCodeNotSpecified(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailExecExitCodeNotSpecified",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)
	debugout := string(outerr)
	ec := 2
	// Yay, different exit codes for the same not found error...
	if runtime.GOOS == "darwin" {
		ec = 1
	}
	msg := fmt.Sprintf(
		"assertion failed: not equal: expected 0 but got %d", ec,
	)
	require.Contains(debugout, msg)
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

func TestFailExecTimeoutPluginDefault(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
	require := require.New(t)

	fp := filepath.Join("testdata", "timeout-plugin-default.yaml")
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
}

func TestExecTimeoutPluginDefault(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailExecTimeoutPluginDefault",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)
	debugout := string(outerr)
	require.Contains(debugout, "using timeout of 0.5s [plugin default]")
	require.Contains(debugout, "assertion failed: timeout exceeded")
}

func TestFailExecSleepTimeout(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
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

func TestExecSleepTimeout(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailExecSleepTimeout",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)
	debugout := string(outerr)
	require.Contains(debugout, "assertion failed: timeout exceeded")
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

func TestFailExecTimeoutCascade(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
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

	ctx := gdtcontext.New(gdtcontext.WithDebug())
	err = s.Run(ctx, t)
	require.Nil(err)
}

func TestExecTimeoutCascade(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailExecTimeoutCascade",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)

	debugout := string(outerr)
	require.Contains(debugout, "using timeout of 500ms [scenario default]")
	require.Contains(debugout, "using timeout of 20ms")
}

func TestFailExecOnFail(t *testing.T) {
	if !*failFlag {
		t.Skip("skipping without -fail flag")
	}
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
}

func TestExecOnFail(t *testing.T) {
	require := require.New(t)
	target := os.Args[0]
	failArgs := []string{
		"-test.v",
		"-test.run=FailExecOnFail",
		"-fail",
	}
	outerr, err := exec.Command(target, failArgs...).CombinedOutput()

	// The test should have failed...
	require.NotNil(err)

	debugout := string(outerr)
	require.Contains(debugout, "assertion failed: not equal: expected dat but got cat")
	require.Contains(debugout, "echo [bad kitty]")
}

func TestTimeoutWithWait(t *testing.T) {
	require := require.New(t)

	fp := filepath.Join("testdata", "timeout-with-wait.yaml")
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
}
