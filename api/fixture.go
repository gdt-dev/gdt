// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import "context"

// A Fixture allows state to be passed from setups
type Fixture interface {
	// Start sets up the fixture
	Start(context.Context) error
	// Stop tears down the fixture, cleaning up any owned resources
	Stop(context.Context)
	// HasState returns true if the fixture contains some state with the given
	// key
	HasState(string) bool
	// State returns the state data at the given key, or nil if no such state
	// key is managed by the fixture
	State(string) interface{}
}
