// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/gdt-dev/gdt/parse"
)

// FromReader parses the supplied io.Reader and returns a Scenario representing
// the contents in the reader. Returns an error if any syntax or validation
// failed
func FromReader(
	r io.Reader,
	mods ...ScenarioModifier,
) (*Scenario, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(contents, mods...)
}

// FromBytes returns a Scenario after parsing the supplied contents
func FromBytes(
	contents []byte,
	mods ...ScenarioModifier,
) (*Scenario, error) {
	s := New(mods...)
	expanded := parse.ExpandWithFixedDoubleDollar(string(contents))
	if err := yaml.Unmarshal([]byte(expanded), s); err != nil {
		return nil, err
	}

	return s, nil
}
