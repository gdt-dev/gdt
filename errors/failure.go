// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrFailure is the base error class for all errors that represent failed
	// assertions when evaluating a test.
	ErrFailure = errors.New("assertion failed")
	// ErrTimeoutExceeded is an ErrFailure when a test's execution exceeds a
	// timeout length.
	ErrTimeoutExceeded = fmt.Errorf("%s: timeout exceeded", ErrFailure)
	// ErrNotEqual is an ErrFailure when an expected thing doesn't equal an
	// observed thing.
	ErrNotEqual = fmt.Errorf("%w: not equal", ErrFailure)
	// ErrNotIn is an ErrFailure when an expected thing doesn't appear in an
	// expected container.
	ErrNotIn = fmt.Errorf("%w: not in", ErrFailure)
	// ErrNoneIn is an ErrFailure when none of a list of elements appears in an
	// expected container.
	ErrNoneIn = fmt.Errorf("%w: none in", ErrFailure)
	// ErrUnexpectedError is an ErrFailure when an unexpected error has
	// occurred.
	ErrUnexpectedError = fmt.Errorf("%w: unexpected error", ErrFailure)
)

// TimeoutExceeded returns an ErrTimeoutExceeded when a test's execution
// exceeds a timeout length. The optional failure parameter indicates a failed
// assertion that occurred before a timeout was reached.
func TimeoutExceeded(duration string, failure error) error {
	if failure != nil {
		return fmt.Errorf(
			"%w: timed out waiting for assertion to succeed (%s)",
			failure, duration,
		)
	}
	return fmt.Errorf("%s (%s)", ErrTimeoutExceeded, duration)
}

// NotEqualLength returns an ErrNotEqual when an expected length doesn't
// equal an observed length.
func NotEqualLength(exp, got int) error {
	return fmt.Errorf(
		"%w: expected length of %d but got %d",
		ErrNotEqual, exp, got,
	)
}

// NotEqual returns an ErrNotEqual when an expected thing doesn't equal an
// observed thing.
func NotEqual(exp, got interface{}) error {
	return fmt.Errorf("%w: expected %v but got %v", ErrNotEqual, exp, got)
}

// NotIn returns an ErrNotIn when an expected thing doesn't appear in an
// expected container.
func NotIn(element, container interface{}) error {
	return fmt.Errorf(
		"%w: expected %v to contain %v",
		ErrNotIn, container, element,
	)
}

// NoneIn returns an ErrNoneIn when none of a list of elements appears in an
// expected container.
func NoneIn(elements, container interface{}) error {
	return fmt.Errorf(
		"%w: expected %v to contain one of %v",
		ErrNoneIn, container, elements,
	)
}

// UnexpectedError returns an ErrUnexpectedError when a supplied error is not
// expected.
func UnexpectedError(err error) error {
	return fmt.Errorf("%w: %s", ErrUnexpectedError, err)
}
