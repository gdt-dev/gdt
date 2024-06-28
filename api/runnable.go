// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"context"
	"testing"
)

// Runnable are things that Run a `*testing.T`
type Runnable interface {
	// Run accepts a context and a `*testing.T` and runs some tests within that
	// context
	//
	// Errors returned by Run() are **RuntimeErrors**, not failures in
	// assertions.
	Run(context.Context, *testing.T) error
}
