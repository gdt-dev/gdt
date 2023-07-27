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
	Run(context.Context, *testing.T) error
}
