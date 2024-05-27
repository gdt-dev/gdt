// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import "context"

// Assertions track zero or more assertions about some result
type Assertions interface {
	// OK returns true if all contained assertions pass successfully, false
	// otherwise. If false is returned, Failures() is guaranteed to be
	// non-empty.
	OK(context.Context) bool
	// Fail appends a supplied error to the set of failed assertions
	Fail(error)
	// Failures returns a slice of failure messages indicating which assertions
	// did not succeed.
	Failures() []error
}
