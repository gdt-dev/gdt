// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/shlex"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/api"
)

var (
	// ErrUnknownShell returns an ErrParse when an unknown shell is specified
	ErrUnknownShell = fmt.Errorf(
		"%w: unknown shell", api.ErrParse,
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
		return api.ExpectedMapAt(node)
	}
	var execValNode *yaml.Node
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "shell":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			s.Shell = strings.TrimSpace(valNode.Value)
			if _, err := exec.LookPath(s.Shell); err != nil {
				return UnknownShell(s.Shell)
			}
		case "exec":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			execValNode = valNode
			s.Exec = strings.TrimSpace(valNode.Value)
			if s.Exec == "" {
				return ExecEmpty(valNode)
			}
		case "assert":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			var e *Expect
			if err := valNode.Decode(&e); err != nil {
				return err
			}
			s.Assert = e
		case "on":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			var o *On
			if err := valNode.Decode(&o); err != nil {
				return err
			}
			s.On = o
		default:
			if lo.Contains(api.BaseSpecFields, key) {
				continue
			}
			return api.UnknownFieldAt(key, keyNode)
		}
	}
	if s.Exec == "" {
		return ExecEmpty(node)
	}
	if s.Shell != "" {
		_, err := shlex.Split(s.Exec)
		if err != nil {
			return ExecInvalidShellParse(err, execValNode)
		}
	}
	return nil
}

func (e *Expect) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return api.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "exit_code", "exit-code":
			if valNode.Kind != yaml.ScalarNode {
				return api.ExpectedScalarAt(valNode)
			}
			ec, err := strconv.Atoi(valNode.Value)
			if err != nil {
				return err
			}
			e.ExitCode = ec
		case "out":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			var pe *PipeExpect
			if err := valNode.Decode(&pe); err != nil {
				return err
			}
			e.Out = pe
		case "err":
			if valNode.Kind != yaml.MappingNode {
				return api.ExpectedMapAt(valNode)
			}
			var pe *PipeExpect
			if err := valNode.Decode(&pe); err != nil {
				return err
			}
			e.Err = pe
		default:
			return api.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}

func (e *PipeExpect) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return api.ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return api.ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "all", "is", "contains", "contains-all", "contains_all":
			if valNode.Kind != yaml.ScalarNode && valNode.Kind != yaml.SequenceNode {
				return api.ExpectedScalarOrSequenceAt(valNode)
			}
			var v api.FlexStrings
			if err := valNode.Decode(&v); err != nil {
				return err
			}
			e.ContainsAll = &v
		case "any", "contains-one-of", "contains-any", "contains_one_of", "contains_any":
			if valNode.Kind != yaml.ScalarNode && valNode.Kind != yaml.SequenceNode {
				return api.ExpectedScalarOrSequenceAt(valNode)
			}
			var v api.FlexStrings
			if err := valNode.Decode(&v); err != nil {
				return err
			}
			e.ContainsAny = &v
		case "none", "none-of", "contains-none-of", "contains-none", "none_of", "contains_none_of", "contains_none":
			if valNode.Kind != yaml.ScalarNode && valNode.Kind != yaml.SequenceNode {
				return api.ExpectedScalarOrSequenceAt(valNode)
			}
			var v api.FlexStrings
			if err := valNode.Decode(&v); err != nil {
				return err
			}
			e.ContainsNone = &v
		default:
			return api.UnknownFieldAt(key, keyNode)
		}
	}
	return nil
}
