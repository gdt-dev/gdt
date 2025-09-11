// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/gdt-dev/gdt/api"
	gdtcontext "github.com/gdt-dev/gdt/context"
	"github.com/gdt-dev/gdt/debug"
	"github.com/google/shlex"
	"github.com/samber/lo"
)

// Action describes a single execution of one or more commands via the
// operating system's `exec` family of functions.
type Action struct {
	// Exec is the exact command to execute.
	//
	// You may execute more than one command but must include the `shell` field
	// to indicate that the command should be run in a shell. It is best
	// practice, however, to simply use multiple `exec` specs instead of
	// executing multiple commands in a single shell call.
	Exec string `yaml:"exec"`
	// Shell is the specific shell to use in executing the command. If empty
	// (the default), no shell is used to execute the command and instead the
	// operating system's `exec` family of calls is used.
	Shell string `yaml:"shell,omitempty"`
	// VarStdout is a shortcut for Var:{VARIABLE_NAME}:from:stdout
	VarStdout string `yaml:"var-stdout,omitempty"`
	// VarStderr is a shortcut for Var:{VARIABLE_NAME}:from:stderr
	VarStderr string `yaml:"var-stderr,omitempty"`
	// VarRC is a shortcut for Var:{VARIABLE_NAME}:from:returncode
	VarRC string `yaml:"var-rc,omitempty"`
}

// Do performs a single command or shell execution returning the corresponding
// exit code and any runtime error. The `outbuf` and `errbuf` buffers will be
// filled with the contents of the command's stdout and stderr pipes
// respectively.
func (a *Action) Do(
	ctx context.Context,
	outbuf *bytes.Buffer,
	errbuf *bytes.Buffer,
	exitcode *int,
) error {
	var target string
	var args []string
	if a.Shell == "" {
		// Parse time already validated exec string parses into valid shell
		// args
		args, _ = shlex.Split(a.Exec)
		target = args[0]
		args = args[1:]
	} else {
		target = a.Shell
		args = []string{"-c", a.Exec}
	}

	target = gdtcontext.ReplaceVariables(ctx, target)
	args = lo.Map(args, func(arg string, _ int) string {
		return gdtcontext.ReplaceVariables(ctx, arg)
	})

	debug.Println(ctx, "exec: %s %s", target, args)

	cmd := exec.CommandContext(ctx, target, args...)

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
		return api.ErrTimeoutExceeded
	}
	if err != nil {
		return err
	}
	if outbuf != nil {
		if _, err = outbuf.ReadFrom(outpipe); err != nil {
			debug.Println(ctx, "exec: error reading from stdout: %s", err)
		}
		if outbuf.Len() > 0 {
			debug.Println(ctx, "exec: stdout: %s", outbuf.String())
		}
	}
	if errbuf != nil {
		if _, err = errbuf.ReadFrom(errpipe); err != nil {
			debug.Println(ctx, "exec: error reading from stderr: %s", err)
		}
		if errbuf.Len() > 0 {
			debug.Println(ctx, "exec: stderr: %s", errbuf.String())
		}
	}

	err = cmd.Wait()
	if gdtcontext.TimedOut(ctx, err) {
		return api.ErrTimeoutExceeded
	}
	if err != nil && exitcode != nil {
		eerr, _ := err.(*exec.ExitError)
		ec := eerr.ExitCode()
		*exitcode = ec
	}
	return nil
}
