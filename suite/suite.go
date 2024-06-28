// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"github.com/gdt-dev/gdt/scenario"
)

// Suite contains zero or more scenarios, one for each YAML file
// representing a Scenario in a given directory
type Suite struct {
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
	// Fixtures specifies an ordered list of fixtures the test suite's test
	// cases depend on.
	Fixtures []string `yaml:"fixtures,omitempty"`
	// Scenarios is a collection of test scenarios in this test suite
	Scenarios []*scenario.Scenario `yaml:"-"`
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

// WithFixtures sets a test suite's Fixtures attribute
func WithFixtures(fixtures []string) SuiteModifier {
	return func(s *Suite) {
		s.Fixtures = fixtures
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

// Append appends a test element to the test suite
func (s *Suite) Append(sc *scenario.Scenario) {
	s.Scenarios = append(s.Scenarios, sc)
}
