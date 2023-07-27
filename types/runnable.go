// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import (
	"context"
	"testing"
)

// Runnable represents things that have a Run() method that accepts a Context
// and a pointer to a testing.T. Example things that implement this interface
// are `gdt.Scenario` and `gdt.Suite`.
type Runnable interface {
	// Run executes the suite or scenario. The error that is returned will
	// always be derived from `gdterrors.RuntimeError` and represents an
	// *unrecoverable* error.
	//
	// Test assertion failures are *not* considered errors. The Scenario.Run()
	// method controls whether `testing.T.Fail()` or `testing.T.Skip()` is
	// called which will mark the test units failed or skipped if a test unit
	// evaluates to false.
	Run(context.Context, *testing.T) error
}
