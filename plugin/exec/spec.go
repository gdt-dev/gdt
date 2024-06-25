// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package exec

import (
	"github.com/gdt-dev/gdt/api"
)

// Spec describes a single Spec that executes one or more commands via the
// operating system's `exec` family of functions.
type Spec struct {
	api.Spec
	Action
	// Assert is an object containing the conditions that the Spec will assert.
	Assert *Expect `yaml:"assert,omitempty"`
	// On is an object containing actions to take upon certain conditions.
	On *On `yaml:"on,omitempty"`
}

func (s *Spec) SetBase(b api.Spec) {
	s.Spec = b
}

func (s *Spec) Base() *api.Spec {
	return &s.Spec
}

func (s *Spec) Retry() *api.Retry {
	return nil
}

func (s *Spec) Timeout() *api.Timeout {
	return nil
}
