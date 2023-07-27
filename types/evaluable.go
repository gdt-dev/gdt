// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import (
	"context"
	"testing"

	"github.com/gdt-dev/gdt/result"
)

// Evaluable represents individual test units in a Scenario
type Evaluable interface {
	// Eval performs an action and evaluates the results of that action,
	// returning a Result that informs the Scenario about what failed or
	// succeeded about the Evaluable's conditions.
	Eval(context.Context, *testing.T) *result.Result
	// SetBase sets the Evaluable's base Spec
	SetBase(Spec)
	// Base returns the Evaluable's base Spec
	Base() *Spec
}
