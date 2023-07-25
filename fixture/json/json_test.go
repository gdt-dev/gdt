// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json_test

import (
	"bytes"
	"testing"

	jsonfix "github.com/gdt-dev/gdt/fixture/json"
	gdttypes "github.com/gdt-dev/gdt/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromString(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	s := `{"book": {"title": "The Cat in the Hat", "year": 1957}}`
	f, err := jsonfix.New(s)

	require.Nil(err)
	require.NotNil(f)
	require.Implements((*gdttypes.Fixture)(nil), f)

	assert.True(f.HasState("$.book.year"))
	assert.Equal("1957", f.State("$.book.year"))
	assert.False(f.HasState("$.book.notexist"))
	assert.Nil(f.State("$.book.notexist"))
}

func TestNewFromBytes(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	b := []byte(`{"book": {"title": "The Cat in the Hat", "year": 1957}}`)
	f, err := jsonfix.New(b)

	require.Nil(err)
	require.NotNil(f)
	require.Implements((*gdttypes.Fixture)(nil), f)

	assert.True(f.HasState("$.book.year"))
	assert.Equal("1957", f.State("$.book.year"))
	assert.False(f.HasState("$.book.notexist"))
	assert.Nil(f.State("$.book.notexist"))
}

func TestNewFromReader(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	b := []byte(`{"book": {"title": "The Cat in the Hat", "year": 1957}}`)
	r := bytes.NewReader(b)
	f, err := jsonfix.New(r)

	require.Nil(err)
	require.NotNil(f)
	require.Implements((*gdttypes.Fixture)(nil), f)

	assert.True(f.HasState("$.book.year"))
	assert.Equal("1957", f.State("$.book.year"))
	assert.False(f.HasState("$.book.notexist"))
	assert.Nil(f.State("$.book.notexist"))
}
