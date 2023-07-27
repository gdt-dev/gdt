// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	gdttypes "github.com/gdt-dev/gdt/types"
)

// Spec describes a single Spec that executes one or more commands via the
// operating system's `exec` family of functions.
type Spec struct {
	gdttypes.Spec
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
	// ExitCode is the expected exit code for the executed command. The default
	// (0) is the universal successful exit code, so you only need to set this
	// if you expect a non-successful result from executing the command.
	ExitCode int `yaml:"exit_code,omitempty"`
	// Out has things that are expected in the stdout response
	Out *PipeExpect `yaml:"out,omitempty"`
	// Err has things that are expected in the stderr response
	Err *PipeExpect `yaml:"err,omitempty"`
}

func (s *Spec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *Spec) Base() *gdttypes.Spec {
	return &s.Spec
}
