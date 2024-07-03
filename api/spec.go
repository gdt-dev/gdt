// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// BaseSpecFields contains the list of base spec fields for plugin Spec
	// types to use in ignoring unknown fields.
	BaseSpecFields = []string{
		"name",
		"description",
		"timeout",
		"wait",
		"retry",
	}
)

// Spec represents a single test action and one or more assertions about
// output or behaviour. All gdt plugins have their own Spec structs that
// inherit from this base struct.
type Spec struct {
	// Plugin is a pointer to the plugin that successfully parsed the test spec
	Plugin Plugin `yaml:"-"`
	// Defaults contains the parsed defaults for the Spec. These are injected
	// by the scenario during parse.
	Defaults *Defaults `yaml:"-"`
	// Index within the scenario where this Spec is located
	Index int `yaml:"-"`
	// Name for the individual test unit
	Name string `yaml:"name,omitempty"`
	// Description of the test unit
	Description string `yaml:"description,omitempty"`
	// Timeout contains the timeout configuration for the Spec
	Timeout *Timeout `yaml:"timeout,omitempty"`
	// Wait contains the wait configuration for the Spec
	Wait *Wait `yaml:"wait,omitempty"`
	// Retry contains the retry configuration for the Spec
	Retry *Retry `yaml:"retry,omitempty"`
}

// Title returns the Name of the scenario or the Path's file/base name if there
// is no name.
func (s *Spec) Title() string {
	if s.Name != "" {
		return s.Name
	}
	if s.Description != "" {
		return slugify(s.Description)
	}
	return strconv.Itoa(s.Index)
}

// slugify returns a new string that lowercases and removes spaces and forward
// slashes from the supplied string
func slugify(s string) string {
	s = strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.TrimSpace(s),
				" ", "-"),
			"/", "-",
		),
	)
	for {
		if strings.Contains(s, "--") {
			s = strings.ReplaceAll(s, "--", "-")
		} else {
			return s
		}
	}
}

// UnmarshalYAML examines the mapping YAML node for base Spec fields and sets
// the associated struct field from that value node.
func (s *Spec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return ExpectedMapAt(node)
	}
	// maps/structs are stored in a top-level Node.Content field which is a
	// concatenated slice of Node pointers in pairs of key/values.
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		if keyNode.Kind != yaml.ScalarNode {
			return ExpectedScalarAt(keyNode)
		}
		key := keyNode.Value
		valNode := node.Content[i+1]
		switch key {
		case "name":
			if valNode.Kind != yaml.ScalarNode {
				return ExpectedScalarAt(valNode)
			}
			s.Name = valNode.Value
		case "description":
			if valNode.Kind != yaml.ScalarNode {
				return ExpectedScalarAt(valNode)
			}
			s.Description = valNode.Value
		case "timeout":
			var to *Timeout
			switch valNode.Kind {
			case yaml.MappingNode:
				// We support the old-style timeout:after
				if err := valNode.Decode(&to); err != nil {
					return ExpectedTimeoutAt(valNode)
				}
			case yaml.ScalarNode:
				// We also support a straight string duration
				to = &Timeout{
					After: valNode.Value,
				}
			default:
				return ExpectedScalarOrMapAt(valNode)
			}
			_, err := time.ParseDuration(to.After)
			if err != nil {
				return err
			}
			s.Timeout = to
		case "wait":
			if valNode.Kind != yaml.MappingNode {
				return ExpectedMapAt(valNode)
			}
			var w *Wait
			if err := valNode.Decode(&w); err != nil {
				return ExpectedWaitAt(valNode)
			}
			if w.Before != "" {
				_, err := time.ParseDuration(w.Before)
				if err != nil {
					return err
				}
			}
			if w.After != "" {
				_, err := time.ParseDuration(w.After)
				if err != nil {
					return err
				}
			}
			s.Wait = w
		case "retry":
			if valNode.Kind != yaml.MappingNode {
				return ExpectedMapAt(valNode)
			}
			var r *Retry
			if err := valNode.Decode(&r); err != nil {
				return ExpectedRetryAt(valNode)
			}
			if r.Attempts != nil {
				attempts := *r.Attempts
				if attempts < 1 {
					return InvalidRetryAttempts(valNode, attempts)
				}
			}
			if r.Interval != "" {
				_, err := time.ParseDuration(r.Interval)
				if err != nil {
					return err
				}
			}
			s.Retry = r
		}
	}
	return nil
}
