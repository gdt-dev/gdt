// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/errors"
	gdttypes "github.com/gdt-dev/gdt/types"
)

var (
	// ErrUnknownShell returns an ErrParse when an unknown shell is specified
	ErrUnknownShell = fmt.Errorf(
		"%w: unknown shell", errors.ErrParse,
	)
)

// UnknownShell returns a wrapped version of ErrParse that indicates the
// user specified an unknown shell.
func UnknownShell(shell string) error {
	return fmt.Errorf(
		"%w: %s", ErrUnknownShell, shell,
	)
}

func (s *Spec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return errors.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return errors.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "shell":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Shell = strings.TrimSpace(valNode.Value)
			if _, err := exec.LookPath(s.Shell); err != nil {
				return UnknownShell(s.Shell)
			}
		case "exec":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			s.Exec = strings.TrimSpace(valNode.Value)
		case "exit_code":
			if valNode.Kind != yaml.ScalarNode {
				return errors.ExpectedScalarAt(valNode)
			}
			ec, err := strconv.Atoi(valNode.Value)
			if err != nil {
				return err
			}
			s.ExitCode = ec
		case "out":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			var pe *PipeExpect
			if err := valNode.Decode(&pe); err != nil {
				return err
			}
			s.Out = pe
		case "err":
			if valNode.Kind != yaml.MappingNode {
				return errors.ExpectedMapAt(valNode)
			}
			var pe *PipeExpect
			if err := valNode.Decode(&pe); err != nil {
				return err
			}
			s.Err = pe
		default:
			if lo.Contains(gdttypes.BaseSpecFields, key) {
				continue
			}
			return errors.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}
