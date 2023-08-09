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
	Action
	// Assert is an object containing the conditions that the Spec will assert.
	Assert *Expect `yaml:"assert,omitempty"`
	// On is an object containing actions to take upon certain conditions.
	On *On `yaml:"on,omitempty"`
}

func (s *Spec) SetBase(b gdttypes.Spec) {
	s.Spec = b
}

func (s *Spec) Base() *gdttypes.Spec {
	return &s.Spec
}
