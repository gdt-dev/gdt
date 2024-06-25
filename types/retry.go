// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package types

import (
	"time"
)

const (
	// DefaultRetryAttempts indicates the default number of times to retry
	// retries when the plugin uses retries but has not specified a number of
	// attempts.
	DefaultRetryAttempts = 3 * time.Second
	// DefaultRetryConstantInterval indicates the default interval to use for
	// retries when the plugin uses retries but does not use exponential
	// backoff.
	DefaultRetryConstantInterval = 3 * time.Second
)

var (
	// NoRetry indicates that there should not be any retry attempts. It is
	// passed from a plugin to indicate a Spec should not be retried.
	NoRetry = &Retry{}
)

// Retry contains information about the number of attempts and interval
// duration with which a Plugin should re-run a Spec's action if the Spec's
// assertions fail.
type Retry struct {
	// Attempts is the number of  times that the test unit should be retried in
	// the event of assertion failure.
	Attempts *int `yaml:"attempts,omitempty"`
	// Interval is the amount of time that the plugin should wait before
	// retrying the test unit in the event of assertion failure.
	// Specify a duration using Go's time duration string.
	// See https://pkg.go.dev/time#ParseDuration
	Interval string `yaml:"interval,omitempty"`
	// Exponential indicates that an exponential backoff should be applied to
	// the retry. When true, the value of Interval, if any, is used as the
	// initial interval for the backoff algoritm.
	Exponential bool `yaml:"exponential,omitempty"`
}

// IntervalDuration returns the time duration of the Retry.Interval
func (r *Retry) IntervalDuration() time.Duration {
	// Parsing already validated the duration string so no need to check again
	// here
	dur, _ := time.ParseDuration(r.Interval)
	return dur
}
