// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"bytes"
	"strings"

	"github.com/gdt-dev/gdt/errors"
	gdttypes "github.com/gdt-dev/gdt/types"
)

// PipeExpect contains assertions about the contents of a pipe
type PipeExpect struct {
	// Is contains the exact match (minus whitespace) of the contents of the
	// pipe
	Is *string `yaml:"is,omitempty"`
	// Contains is one or more strings that *all* must be present in the
	// contents of the pipe
	Contains []string `yaml:"contains,omitempty"`
	// ContainsOneOf is one or more strings of which *at least one* must be
	// present in the contents of the pipe
	ContainsOneOf []string `yaml:"contains_one_of,omitempty"`
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
	// terminal indicates there was a failure in evaluating the assertions that
	// should be considered a terminal condition (and therefore the test action
	// should not be retried).
	terminal bool
}

// Fail appends a supplied error to the set of failed assertions
func (a *pipeAssertions) Fail(err error) {
	a.failures = append(a.failures, err)
}

// Failures returns a slice of errors for all failed assertions
func (a *pipeAssertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// Terminal returns a bool indicating the assertions failed in a way that is
// not retryable.
func (a *pipeAssertions) Terminal() bool {
	if a == nil {
		return false
	}
	return a.terminal
}

// OK checks all the assertions in the pipeAssertions against the supplied pipe
// contents and returns true if all assertions pass.
func (a *pipeAssertions) OK() bool {
	if a == nil || a.pipe == nil {
		return true
	}

	res := true
	contents := strings.TrimSpace(a.pipe.String())
	if a.Is != nil {
		exp := *a.Is
		got := contents
		if exp != got {
			a.Fail(errors.NotEqual(exp, got))
			res = false
		}
	}
	if len(a.Contains) > 0 {
		for _, find := range a.Contains {
			if !strings.Contains(contents, find) {
				a.Fail(errors.NotIn(find, a.name))
				res = false
			}
		}
	}
	if len(a.ContainsOneOf) > 0 {
		found := false
		for _, find := range a.ContainsOneOf {
			if idx := strings.Index(contents, find); idx > -1 {
				found = true
				break
			}
		}
		if !found {
			a.Fail(errors.NoneIn(a.ContainsOneOf, a.name))
			res = false
		}
	}
	return res
}

// assertions contains all assertions made for the exec test
type assertions struct {
	// failures contains the set of error messages for failed assertions
	failures []error
	// terminal indicates there was a failure in evaluating the assertions that
	// should be considered a terminal condition (and therefore the test action
	// should not be retried).
	terminal bool
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

// Failures returns a slice of errors for all failed assertions
func (a *assertions) Failures() []error {
	if a == nil {
		return []error{}
	}
	return a.failures
}

// Terminal returns a bool indicating the assertions failed in a way that is
// not retryable.
func (a *assertions) Terminal() bool {
	if a == nil {
		return false
	}
	return a.terminal
}

// OK checks all the assertions against the supplied arguments and returns true
// if all assertions pass.
func (a *assertions) OK() bool {
	res := true
	if a.expExitCode != a.exitCode {
		a.Fail(errors.NotEqual(a.expExitCode, a.exitCode))
		res = false
	}
	if !a.expOutPipe.OK() {
		a.failures = append(a.failures, a.expOutPipe.Failures()...)
		res = false
	}
	if !a.expErrPipe.OK() {
		a.failures = append(a.failures, a.expErrPipe.Failures()...)
		res = false
	}
	return res
}

// newAssertions returns an assertions object populated with the supplied exec
// spec assertions
func newAssertions(
	expExitCode int,
	exitCode int,
	expOutPipe *PipeExpect,
	outPipe *bytes.Buffer,
	expErrPipe *PipeExpect,
	errPipe *bytes.Buffer,
) gdttypes.Assertions {
	a := &assertions{
		failures:    []error{},
		expExitCode: exitCode,
		exitCode:    exitCode,
	}
	if expOutPipe != nil {
		a.expOutPipe = &pipeAssertions{
			PipeExpect: *expOutPipe,
			name:       "stdout",
			pipe:       outPipe,
		}
	}
	if expErrPipe != nil {
		a.expErrPipe = &pipeAssertions{
			PipeExpect: *expErrPipe,
			name:       "stderr",
			pipe:       errPipe,
		}
	}
	return a
}
