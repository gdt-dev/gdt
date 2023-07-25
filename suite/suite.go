// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"context"

	gdttypes "github.com/gdt-dev/gdt/types"
)

// Suite contains zero or more Runnable things, one for each YAML file
// representing a Scenario in a given directory
type Suite struct {
	// ctx stores the context. Yes, I know this is not good practice and that a
	// context should be passed as the first argument to all methods, but the
	// `yaml.Unmarshaler` interface does not have a context argument and
	// there's no other way to pass in necessary information.
	ctx context.Context
	// Path is the filepath to the test suite directory.
	Path string `yaml:"-"`
	// Name is the short name for the test suite. If empty, defaults to Path.
	Name string `yaml:"name,omitempty"`
	// Description is a description of the tests contained in the test suite.
	Description string `yaml:"description,omitempty"`
	// Defaults contains any default configuration values for test cases
	// contained within the test suite.
	//
	// During parsing, plugins are handed this raw data and asked to interpret
	// it into known configuration values for that plugin.
	Defaults map[string]interface{} `yaml:"defaults,omitempty"`
	// Require specifies an ordered list of fixtures the test suite's test
	// cases depends on.
	Require []string `yaml:"require,omitempty"`
	// Scenarios is a collection of test scenarios in this test suite
	Scenarios []gdttypes.Runnable `yaml:"-"`
}

// SuiteModifier sets some value on the test suite
type SuiteModifier func(s *Suite)

// WithName sets a test suite's Name attribute
func WithName(name string) SuiteModifier {
	return func(s *Suite) {
		s.Name = name
	}
}

// WithPath sets a test suite's Path attribute
func WithPath(path string) SuiteModifier {
	return func(s *Suite) {
		s.Path = path
	}
}

// WithDescription sets a test suite's Description attribute
func WithDescription(description string) SuiteModifier {
	return func(s *Suite) {
		s.Description = description
	}
}

// WithDefaults sets a test suite's Defaults attribute
func WithDefaults(defaults map[string]interface{}) SuiteModifier {
	return func(s *Suite) {
		s.Defaults = defaults
	}
}

// WithRequires sets a test suite's Requires attribute
func WithRequires(require []string) SuiteModifier {
	return func(s *Suite) {
		s.Require = require
	}
}

// WithContext sets a test scenario's context
func WithContext(ctx context.Context) SuiteModifier {
	return func(s *Suite) {
		s.ctx = ctx
	}
}

// New returns a new Suite
func New(mods ...SuiteModifier) *Suite {
	s := &Suite{}
	for _, mod := range mods {
		mod(s)
	}
	return s
}

// Append appends a runnable test element to the test suite
func (s *Suite) Append(r gdttypes.Runnable) {
	s.Scenarios = append(s.Scenarios, r)
}
