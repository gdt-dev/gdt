// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	gopath "path"

	gdttypes "github.com/gdt-dev/gdt/types"
)

// Scenario is a generalized gdt test case file. It contains a set of Runnable
// test units.
type Scenario struct {
	// Path is the filepath to the test case.
	Path string `yaml:"-"`
	// Name is the short name for the test case. If empty, defaults to the base
	// filename in Path.
	Name string `yaml:"name,omitempty"`
	// Description is a description of the tests contained in the test case.
	Description string `yaml:"description,omitempty"`
	// Defaults contains any default configuration values for test specs
	// contained within the test scenario.
	//
	// During parsing, plugins are handed this raw data and asked to interpret
	// it into known configuration values for that plugin.
	Defaults map[string]interface{} `yaml:"defaults,omitempty"`
	// Fixtures specifies an ordered list of fixtures the test case depends on.
	Fixtures []string `yaml:"fixtures,omitempty"`
	// Tests is the collection of test units in this test case. These will be
	// the fully parsed and materialized plugin Spec structs.
	Tests []gdttypes.TestUnit `yaml:"tests,omitempty"`
}

// Title returns the Name of the scenario or the Path's file/base name if there
// is no name.
func (s *Scenario) Title() string {
	if s.Name != "" {
		return s.Name
	}
	return gopath.Base(s.Path)
}

// ScenarioModifier sets some value on the test scenario
type ScenarioModifier func(s *Scenario)

// WithName sets a test scenario's Name attribute
func WithName(name string) ScenarioModifier {
	return func(s *Scenario) {
		s.Name = name
	}
}

// WithPath sets a test scenario's Path attribute
func WithPath(path string) ScenarioModifier {
	return func(s *Scenario) {
		s.Path = path
	}
}

// WithDescription sets a test scenario's Description attribute
func WithDescription(description string) ScenarioModifier {
	return func(s *Scenario) {
		s.Description = description
	}
}

// WithDefaults sets a test scenario's Defaults attribute
func WithDefaults(defaults map[string]interface{}) ScenarioModifier {
	return func(s *Scenario) {
		s.Defaults = defaults
	}
}

// WithFixtures sets a test scenario's Fixtures attribute
func WithRequires(fixtures []string) ScenarioModifier {
	return func(s *Scenario) {
		s.Fixtures = fixtures
	}
}

// New returns a new Scenario
func New(mods ...ScenarioModifier) *Scenario {
	s := &Scenario{
		Defaults: map[string]interface{}{},
	}
	for _, mod := range mods {
		mod(s)
	}
	return s
}
