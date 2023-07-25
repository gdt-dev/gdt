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
)

// Run executes the specific exec test spec.
func (s *Spec) Run(ctx context.Context, t *testing.T) error {
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}

	var err error
	var cmd *exec.Cmd
	var target string
	var args []string
	if s.Shell == "" {
		args, err = shlex.Split(s.Exec)
		if err != nil {
			return err
		}
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
		return err
	}
	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if gdtcontext.TimedOut(ctx, err) {
		return gdterrors.ErrTimeout
	}
	if err != nil {
		return err
	}
	outbuf.ReadFrom(outpipe)
	errbuf.ReadFrom(errpipe)

	err = cmd.Wait()
	if gdtcontext.TimedOut(ctx, err) {
		return gdterrors.ErrTimeout
	}
	ec := 0
	if err != nil {
		eerr, _ := err.(*exec.ExitError)
		ec = eerr.ExitCode()
	}
	assertions := newAssertions(
		s.ExitCode, ec, s.Out, outbuf, s.Err, errbuf,
	)

	if !assertions.OK() {
		for _, failure := range assertions.Failures() {
			t.Error(failure)
		}
	}
	return nil
}
