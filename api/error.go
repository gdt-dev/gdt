// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
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
	// ErrIn is an ErrFailure when a thing unexpectedly appears in an
	// container.
	ErrIn = fmt.Errorf("%w: in", ErrFailure)
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

// In returns an ErrIn when a thing unexpectedly appears in a container.
func In(element, container interface{}) error {
	return fmt.Errorf(
		"%w: expected %v not to contain %v",
		ErrIn, container, element,
	)
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

var (
	// ErrUnknownSourceType indicates that a From() function was called with an
	// unknown source parameter type.
	ErrUnknownSourceType = errors.New("unknown source argument type")
	// ErrUnknownSpec indicates that there was a test spec definition in a YAML
	// file that no plugin could parse.
	ErrUnknownSpec = errors.New("no plugin could parse spec definition")
	// ErrUnknownField indicates that there was an unknown field in the parsing
	// of a spec or scenario.
	ErrUnknownField = errors.New("unknown field")
	// ErrParse indicates a YAML definition is not valid
	ErrParse = errors.New("invalid YAML")
	// ErrExpectedMap indicates that we did not find an expected mapping
	// field
	ErrExpectedMap = fmt.Errorf(
		"%w: expected map field", ErrParse,
	)
	// ErrExpectedScalar indicates that we did not find an expected scalar
	// field
	ErrExpectedScalar = fmt.Errorf(
		"%w: expected scalar field", ErrParse,
	)
	// ErrExpectedSequence indicates that we did not find an expected
	// scalar field
	ErrExpectedSequence = fmt.Errorf(
		"%w: expected sequence field", ErrParse,
	)
	// ErrExpectedInt indicates that we did not find an expected integer
	// value
	ErrExpectedInt = fmt.Errorf(
		"%w: expected int value", ErrParse,
	)
	// ErrExpectedScalarOrMap indicates that we did not find an expected
	// scalar or map field
	ErrExpectedScalarOrMap = fmt.Errorf(
		"%w: expected scalar or map field", ErrParse,
	)
	// ErrExpectedScalarOrSequence indicates that we did not find an expected
	// scalar or sequence of scalars field
	ErrExpectedScalarOrSequence = fmt.Errorf(
		"%w: expected scalar or sequence of scalars field", ErrParse,
	)
	// ErrExpectedTimeout indicates that the timeout specification was not
	// valid.
	ErrExpectedTimeout = fmt.Errorf(
		"%w: expected timeout specification", ErrParse,
	)
	// ErrExpectedWait indicates that the wait specification was not valid.
	ErrExpectedWait = fmt.Errorf(
		"%w: expected wait specification", ErrParse,
	)
	// ErrExpectedRetry indicates that the retry specification was not valid.
	ErrExpectedRetry = fmt.Errorf(
		"%w: expected retry specification", ErrParse,
	)
	// ErrInvalidRetryAttempts indicates that the retry attempts was not
	// positive.
	ErrInvalidRetryAttempts = fmt.Errorf(
		"%w: invalid retry attempts", ErrParse,
	)
	// ErrFileNotFound is returned when a file path does not exist for a
	// create/apply/delete target.
	ErrFileNotFound = fmt.Errorf(
		"%w: file not found", ErrParse,
	)
)

// UnknownSpecAt returns an ErrUnknownSpec with the line/column of the supplied
// YAML node.
func UnknownSpecAt(path string, node *yaml.Node) error {
	return fmt.Errorf(
		"%w in %s at line %d, column %d",
		ErrUnknownSpec, path, node.Line, node.Column,
	)
}

// UnknownFieldAt returns an ErrUnknownField for a supplied field annotated
// with the line/column of the supplied YAML node.
func UnknownFieldAt(field string, node *yaml.Node) error {
	return fmt.Errorf(
		"%w: %q at line %d, column %d",
		ErrUnknownField, field, node.Line, node.Column,
	)
}

// ExpectedMapAt returns an ErrExpectedMap error annotated with the
// line/column of the supplied YAML node.
func ExpectedMapAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedMap, node.Line, node.Column,
	)
}

// ExpectedScalarAt returns an ErrExpectedScalar error annotated with
// the line/column of the supplied YAML node.
func ExpectedScalarAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedScalar, node.Line, node.Column,
	)
}

// ExpectedSequenceAt returns an ErrExpectedSequence error annotated
// with the line/column of the supplied YAML node.
func ExpectedSequenceAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedSequence, node.Line, node.Column,
	)
}

// ExpectedIntAt returns an ErrExpectedInt error annotated
// with the line/column of the supplied YAML node.
func ExpectedIntAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedInt, node.Line, node.Column,
	)
}

// ExpectedScalarOrSequenceAt returns an ErrExpectedScalarOrSequence error
// annotated with the line/column of the supplied YAML node.
func ExpectedScalarOrSequenceAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedScalarOrSequence, node.Line, node.Column,
	)
}

// ExpectedScalarOrMapAt returns an ErrExpectedScalarOrMap error annotated with
// the line/column of the supplied YAML node.
func ExpectedScalarOrMapAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedScalarOrMap, node.Line, node.Column,
	)
}

// ExpectedTimeoutAt returns an ErrExpectedTimeout error annotated
// with the line/column of the supplied YAML node.
func ExpectedTimeoutAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedTimeout, node.Line, node.Column,
	)
}

// ExpectedWaitAt returns an ErrExpectedWait error annotated with the
// line/column of the supplied YAML node.
func ExpectedWaitAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedWait, node.Line, node.Column,
	)
}

// ExpectedRetryAt returns an ErrExpectedRetry error annotated with the
// line/column of the supplied YAML node.
func ExpectedRetryAt(node *yaml.Node) error {
	return fmt.Errorf(
		"%w at line %d, column %d",
		ErrExpectedRetry, node.Line, node.Column,
	)
}

// InvalidRetryAttempts returns an ErrInvalidRetryAttempts error annotated with
// the line/column of the supplied YAML node.
func InvalidRetryAttempts(node *yaml.Node, attempts int) error {
	return fmt.Errorf(
		"%w of %d at line %d, column %d",
		ErrInvalidRetryAttempts, attempts, node.Line, node.Column,
	)
}

// UnknownSourceType returns an ErrUnknownSourceType error describing the
// supplied parameter type.
func UnknownSourceType(source interface{}) error {
	return fmt.Errorf("%w: %T", ErrUnknownSourceType, source)
}

// FileNotFound returns ErrFileNotFound for a given file path
func FileNotFound(path string, node *yaml.Node) error {
	return fmt.Errorf(
		"%w: %s at line %d, column %d",
		ErrFileNotFound, path, node.Line, node.Column,
	)
}

var (
	// RuntimeError is the base error class for all errors occurring during
	// runtime (and not during the parsing of a scenario or spec)
	RuntimeError = errors.New("runtime error")
	// ErrRequiredFixture is returned when a required fixture has not
	// been registered with the context.
	ErrRequiredFixture = fmt.Errorf(
		"%w: required fixture missing",
		RuntimeError,
	)
)

// RequiredFixtureMissing returns an ErrRequiredFixture with the supplied
// fixture name
func RequiredFixtureMissing(name string) error {
	return fmt.Errorf("%w: %s", ErrRequiredFixture, name)
}
