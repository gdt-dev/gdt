// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"time"
)

type SetOn int

const (
	SetOnNone          SetOn = iota
	SetOnSpec                // a test spec override
	SetOnPlugin              // a plugin override
	SetOnPluginDefault       // a plugin default
	SetOnDefault             // a scenario default
)

// Timings contains information about a test scenario's maximum wait and
// timeout duration and what aspect of the scenario (the scenario defaults, a
// plugin default, a test spec override, etc) had the maximum timeout or wait
// value.
//
// We use this information when initially assessing whether the Go test tool's
// overall timeout is shorter than this maximum in order to inform the user to
// increase the Go test tool timeout.
type Timings struct {
	// GoTestTimeout will be the duration of the timeout specified (or
	// defaulted) by the Go test tool
	GoTestTimeout time.Duration
	// TotalWait will be non-zero when there is a wait specified for either the
	// scenario or a test spec and will contain the aggregate duration of all
	// waits
	TotalWait time.Duration
	// MaxTimeout will be non-zero when there is a timeout specified for either
	// the scenario or a test spec and will contain the duration of the maximum
	// timeout
	MaxTimeout time.Duration
	//TimeoutSetOn indicates where the MaxTimeout value was found
	MaxTimeoutSetOn SetOn
	// TimeoutSpecIndex indicates the test spec's index within the scenario where
	// the max timeout was found
	MaxTimeoutSpecIndex int
}

// AddWait adds a wait duration to the Timings and (re)-calculates the Timings'
// MaxWait attributes
func (t *Timings) AddWait(
	d time.Duration,
) {
	t.TotalWait += d
}

// AddTimeout adds a timeout duration to the Timings and (re)-calculates the
// Timings' MaxTimeout attributes
func (t *Timings) AddTimeout(
	d time.Duration,
	on SetOn,
	specIndex int,
) {
	if d == 0 {
		return
	}
	if d.Abs() > t.MaxTimeout {
		t.MaxTimeout = d
		t.MaxTimeoutSetOn = on
		t.MaxTimeoutSpecIndex = specIndex
	}
}
