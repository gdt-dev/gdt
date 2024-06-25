// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api_test

import (
	"testing"

	"github.com/gdt-dev/gdt/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type foo struct {
	Foo api.FlexStrings `yaml:"foo"`
}

type foop struct {
	Foo *api.FlexStrings `yaml:"foo"`
}

func TestFlexStringsError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var f foo
	contents := []byte(`foo: {bar: {baz: 123}}`)
	err := yaml.Unmarshal(contents, &f)

	require.NotNil(err)
	assert.ErrorIs(err, api.ErrExpectedScalarOrSequence)
}

func TestFlexStrings(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var f foo
	contents := []byte(`foo: singlestring`)
	err := yaml.Unmarshal(contents, &f)

	require.Nil(err)
	assert.Equal([]string{"singlestring"}, f.Foo.Values())

	contents = []byte(`foo: [one, two]`)
	err = yaml.Unmarshal(contents, &f)

	require.Nil(err)
	assert.Equal([]string{"one", "two"}, f.Foo.Values())
}

func TestFlexStringsPointer(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var f foop
	contents := []byte(`foo: singlestring`)
	err := yaml.Unmarshal(contents, &f)

	require.Nil(err)
	assert.Equal([]string{"singlestring"}, f.Foo.Values())

	contents = []byte(`foo: [one, two]`)
	err = yaml.Unmarshal(contents, &f)

	require.Nil(err)
	assert.Equal([]string{"one", "two"}, f.Foo.Values())
}
