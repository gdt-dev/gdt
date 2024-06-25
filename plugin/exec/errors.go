// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
)

var (
	// ErrExecEmpty indicates that the user specified an empty "exec"
	// field
	ErrExecEmpty = fmt.Errorf(
		"%w: expected non-empty exec field", api.ErrParse,
	)
	// ErrExecInvalid indicates that the user specified an invalid "exec" field
	ErrExecInvalid = fmt.Errorf(
		"%w: invalid exec field", api.ErrParse,
	)
)

// ExecEmpty returns an ErrExecEmpty with the line/column of the supplied YAML
// node.
func ExecEmpty(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExecEmpty, node.Line, node.Column,
	)
}

// ExecInvalidShellParse returns an ErrExecInvalid with the error from
// shlex.Split
func ExecInvalidShellParse(err error, node *yaml.Node) error {
	return fmt.Errorf(
		"%w: cannot parse shell args: %s at line %d, column %d",
		ErrExecInvalid, err, node.Line, node.Column,
	)
}

// ExecRuntimeError returns a RuntimeError with an error from the Exec() call.
func ExecRuntimeError(err error) error {
	return fmt.Errorf("%w: %s", api.RuntimeError, err)
}
