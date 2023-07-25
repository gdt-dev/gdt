// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package parse_test

import (
	"testing"

	"github.com/gdt-dev/gdt/parse"
	"github.com/stretchr/testify/assert"
)

func TestExpandWithFixedDoubleDollar(t *testing.T) {
	assert := assert.New(t)

	t.Setenv("foo", "bar")
	t.Setenv("bar", "baz")

	cases := []struct {
		content string
		exp     string
	}{
		{
			content: `This is content with no env var expansion`,
			exp:     `This is content with no env var expansion`,
		},
		{
			content: `This is content with $foo`,
			exp:     `This is content with bar`,
		},
		{
			content: `This is content with ${foo}`,
			exp:     `This is content with bar`,
		},
		{
			content: `This is content with $unknown`,
			exp:     `This is content with `,
		},
		{
			content: `This is content with $$LOCATION`,
			exp:     `This is content with $LOCATION`,
		},
	}
	for _, c := range cases {
		got := parse.ExpandWithFixedDoubleDollar(c.content)
		assert.Equal(c.exp, got)
	}
}
