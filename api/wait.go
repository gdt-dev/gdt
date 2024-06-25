// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"time"
)

// Wait contains information about the duration within which a Spec should
// run along with whether a deadline exceeded/timeout error should be expected
// or not.
type Wait struct {
	// Before is the amount of time that the test unit should wait before
	// executing its action.
	// Specify a duration using Go's time duration string.
	// See https://pkg.go.dev/time#ParseDuration
	Before string `yaml:"before,omitempty"`
	// After is the amount of time that the test unit should wait after
	// executing its action.
	// Specify a duration using Go's time duration string.
	// See https://pkg.go.dev/time#ParseDuration
	After string `yaml:"after,omitempty"`
}

// BeforeDuration returns the time duration of the Wait.Before
func (w *Wait) BeforeDuration() time.Duration {
	// Parsing already validated the duration string so no need to check again
	// here
	dur, _ := time.ParseDuration(w.Before)
	return dur
}

// AfterDuration returns the time duration of the Wait.After
func (w *Wait) AfterDuration() time.Duration {
	// Parsing already validated the duration string so no need to check again
	// here
	dur, _ := time.ParseDuration(w.After)
	return dur
}
