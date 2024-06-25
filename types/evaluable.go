// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import (
	"context"

	"github.com/gdt-dev/gdt/result"
)

// Evaluable represents individual test units in a Scenario
type Evaluable interface {
	// Eval performs an action and evaluates the results of that action,
	// returning a Result that informs the Scenario about what failed or
	// succeeded about the Evaluable's conditions.
	//
	// Errors returned by Eval() are **RuntimeErrors**, not failures in
	// assertions.
	Eval(context.Context) (*result.Result, error)
	// SetBase sets the Evaluable's base Spec
	SetBase(Spec)
	// Base returns the Evaluable's base Spec
	Base() *Spec
	// Retry returns the Evaluable's Retry override, if any
	Retry() *Retry
	// Timeout returns the Evaluable's Timeout override, if any
	Timeout() *Timeout
}
