// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

// TestUnit represents individual test units in a Scenario
type TestUnit interface {
	Runnable
	// SetBase sets the TestUnit's base Spec
	SetBase(Spec)
	// Base returns the TestUnit's base Spec
	Base() *Spec
}
