// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/google/shlex"

	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/gdt-dev/gdt/result"
)

// Eval performs an action and evaluates the results of that action, returning
// a Result that informs the Scenario about what failed or succeeded about the
// Evaluable's conditions.
func (s *Spec) Eval(ctx context.Context, t *testing.T) *result.Result {
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}

	var err error
	var cmd *exec.Cmd
	var target string
	var args []string
	if s.Shell == "" {
		// Parse time already validated exec string parses into valid shell
		// args
		args, _ = shlex.Split(s.Exec)
		target = args[0]
		args = args[1:]
	} else {
		target = s.Shell
		args = []string{"-c", s.Exec}
	}

	debug.Println(ctx, t, "exec: %s %s", target, args)
	cmd = exec.CommandContext(ctx, target, args...)

	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return result.New(result.WithRuntimeError(ExecRuntimeError(err)))
	}
	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return result.New(result.WithRuntimeError(ExecRuntimeError(err)))
	}

	err = cmd.Start()
	if gdtcontext.TimedOut(ctx, err) {
		return result.New(result.WithFailures(gdterrors.ErrTimeoutExceeded))
	}
	if err != nil {
		return result.New(result.WithRuntimeError(ExecRuntimeError(err)))
	}
	outbuf.ReadFrom(outpipe)
	errbuf.ReadFrom(errpipe)

	err = cmd.Wait()
	if gdtcontext.TimedOut(ctx, err) {
		return result.New(result.WithFailures(gdterrors.ErrTimeoutExceeded))
	}
	ec := 0
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		ec = eerr.ExitCode()
	}
	assertions := newAssertions(
		s.ExitCode, ec, s.Out, outbuf, s.Err, errbuf,
	)
	return result.New(result.WithFailures(assertions.Failures()...))
}
