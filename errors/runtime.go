// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors

import (
	"errors"
	"fmt"
)

var (
	// RuntimeError is the base error class for all errors occurring during
	// runtime (and not during the parsing of a scenario or spec)
	RuntimeError = errors.New("runtime error")
	// ErrRequiredFixture is returned when a required fixture has not
	// been registered with the context.
	ErrRequiredFixture = fmt.Errorf(
		"%w: required fixture missing",
		RuntimeError,
	)
)

// RequiredFixtureMissing returns an ErrRequiredFixture with the supplied
// fixture name
func RequiredFixtureMissing(name string) error {
	return fmt.Errorf("%w: %s", ErrRequiredFixture, name)
}
