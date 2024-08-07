// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"context"
	"strings"

	"github.com/gdt-dev/gdt/api"
)

// Expect contains the assertions about an Exec Spec's actions
type Expect struct {
	// ExitCode is the expected exit code for the executed command. The default
	// (0) is the universal successful exit code, so you only need to set this
	// if you expect a non-successful result from executing the command.
	ExitCode int `yaml:"exit-code,omitempty"`
	// Out has things that are expected in the stdout response
	Out *PipeExpect `yaml:"out,omitempty"`
	// Err has things that are expected in the stderr response
	Err *PipeExpect `yaml:"err,omitempty"`
}

// PipeExpect contains assertions about the contents of a pipe
type PipeExpect struct {
	// ContainsAll is one or more strings that *all* must be present in the
	// contents of the pipe
	ContainsAll *api.FlexStrings `yaml:"contains,omitempty"`
	// ContainsNone is one or more strings, *none of which* should be present in
	// the contents of the pipe
	ContainsNone *api.FlexStrings `yaml:"contains-none-of,omitempty"`
	// ContainsOneOf is one or more strings of which *at least one* must be
	// present in the contents of the pipe
	ContainsAny *api.FlexStrings `yaml:"contains-one-of,omitempty"`
}

// pipeAssertions contains assertions about the contents of a pipe
type pipeAssertions struct {
	PipeExpect
	// pipe is the contents of the pipe that we will evaluate.
	pipe *bytes.Buffer
	// name is the string name of the pipe.
	name string
	// failures contains the set of error messages for failed assertions.
	failures []error
}

// Fail appends a supplied error to the set of failed assertions
func (a *pipeAssertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of api for all failed assertions
func (a *pipeAssertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// OK checks all the assertions in the pipeAssertions against the supplied pipe
// contents and returns true if all assertions pass.
func (a *pipeAssertions) OK(ctx context.Context) bool {
	if a == nil || a.pipe == nil {
		return true
	}

	res := true
	contents := strings.TrimSpace(a.pipe.String())
	if a.ContainsAll != nil {
		// When there is just a single value, we use the NotEqual error,
		// otherwise we use the NotIn error
		vals := a.ContainsAll.Values()
		if len(vals) == 1 {
			if !strings.Contains(contents, vals[0]) {
				a.Fail(api.NotEqual(vals[0], contents))
				res = false
			}
		} else {
			for _, find := range vals {
				if !strings.Contains(contents, find) {
					a.Fail(api.NotIn(find, a.name))
					res = false
				}
			}
		}
	}
	if a.ContainsAny != nil {
		found := false
		for _, find := range a.ContainsAny.Values() {
			if idx := strings.Index(contents, find); idx > -1 {
				found = true
				break
			}
		}
		if !found {
			a.Fail(api.NoneIn(a.ContainsAny.Values(), a.name))
			res = false
		}
	}
	if a.ContainsNone != nil {
		for _, find := range a.ContainsNone.Values() {
			if strings.Contains(contents, find) {
				a.Fail(api.In(find, a.name))
				res = false
			}
		}
	}
	return res
}

// assertions contains all assertions made for the exec test
type assertions struct {
	// failures contains the set of error messages for failed assertions
	failures []error
	// expExitCode contains the expected exit code
	expExitCode int
	// exitCode is the exit code we got from the execution
	exitCode int
	// expOutPipe contains the assertions against stdout
	expOutPipe *pipeAssertions
	// expErrPipe contains the assertions against stderr
	expErrPipe *pipeAssertions
}

// Fail appends a supplied error to the set of failed assertions
func (a *assertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of api for all failed assertions
func (a *assertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// OK checks all the assertions against the supplied arguments and returns true
// if all assertions pass.
func (a *assertions) OK(ctx context.Context) bool {
	res := true
	if a.expExitCode != a.exitCode {
		a.Fail(api.NotEqual(a.expExitCode, a.exitCode))
		res = false
	}
	if !a.expOutPipe.OK(ctx) {
		a.failures = append(a.failures, a.expOutPipe.Failures()...)
		res = false
	}
	if !a.expErrPipe.OK(ctx) {
		a.failures = append(a.failures, a.expErrPipe.Failures()...)
		res = false
	}
	return res
}

// newAssertions returns an assertions object populated with the supplied exec
// spec assertions
func newAssertions(
	e *Expect,
	exitCode int,
	outPipe *bytes.Buffer,
	errPipe *bytes.Buffer,
) api.Assertions {
	expExitCode := 0
	if e != nil {
		expExitCode = e.ExitCode
	}
	a := &assertions{
		failures:    []error{},
		expExitCode: expExitCode,
		exitCode:    exitCode,
	}
	if e != nil {
		if e.Out != nil {
			a.expOutPipe = &pipeAssertions{
				PipeExpect: *e.Out,
				name:       "stdout",
				pipe:       outPipe,
			}
		}
		if e.Err != nil {
			a.expErrPipe = &pipeAssertions{
				PipeExpect: *e.Err,
				name:       "stderr",
				pipe:       errPipe,
			}
		}
	}
	return a
}
