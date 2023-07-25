// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package scenario

import (
	"io"
	"io/ioutil"

	"github.com/gdt-dev/gdt/parse"
	gdttypes "github.com/gdt-dev/gdt/types"
	"gopkg.in/yaml.v3"
)

// FromReader parses the supplied io.Reader and returns a Scenario representing
// the contents in the reader. Returns an error if any syntax or validation
// failed
func FromReader(
	r io.Reader,
	mods ...ScenarioModifier,
) (gdttypes.Runnable, error) {
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
) (gdttypes.Runnable, error) {
	s := New(mods...)
	expanded := parse.ExpandWithFixedDoubleDollar(string(contents))
	if err := yaml.Unmarshal([]byte(expanded), s); err != nil {
		return nil, err
	}

	return s, nil
}
