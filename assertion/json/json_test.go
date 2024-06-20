// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gdtjson "github.com/gdt-dev/gdt/assertion/json"
	gdterrors "github.com/gdt-dev/gdt/errors"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestUnsupportedJSONSchemaReference(t *testing.T) {
	require := require.New(t)

	var exp gdtjson.Expect

	// http lookups are not allowed...
	content := []byte(`
schema: http://example.com/schema
`)
	err := yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdtjson.ErrUnsupportedJSONSchemaReference)
}

func TestJSONSchemaFileNotFound(t *testing.T) {
	require := require.New(t)

	var exp gdtjson.Expect

	content := []byte(`
schema: file:///path/does/not/exist
`)
	err := yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdtjson.ErrJSONSchemaFileNotFound)
}

func TestJSONPathInvalid(t *testing.T) {
	require := require.New(t)

	var exp gdtjson.Expect

	content := []byte(`
len: foo
`)
	err := yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorContains(err, "yaml: unmarshal errors")

	content = []byte(`
len: 1
paths: notamap
`)
	err = yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdterrors.ErrExpectedMap)

	content = []byte(`
len: 1
paths:
  noroot: value
`)
	err = yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdtjson.ErrJSONPathInvalidNoRoot)

	content = []byte(`
len: 1
paths:
  $[1-2,3].key: value
`)
	err = yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdtjson.ErrJSONPathInvalid)

	content = []byte(`
len: 1
paths:
  $.: value
`)
	err = yaml.Unmarshal(content, &exp)
	require.NotNil(err)
	require.ErrorIs(err, gdtjson.ErrJSONPathInvalid)
}

func content() []byte {
	b, _ := os.ReadFile(filepath.Join("testdata", "books.json"))
	return b
}

func TestLength(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()
	expLen := len(c)

	exp := gdtjson.Expect{
		Len: &expLen,
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK(ctx))
	require.Empty(a.Failures())

	expLen = 0
	a = gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdterrors.ErrNotEqual)
}

func TestJSONUnmarshalError(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := []byte(`not { value } json`)

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"1234": "foo",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONUnmarshalError)
}

func TestJSONPathError(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"[0].pages": "127",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathNotFound)
}

func TestJSONPathConversionError(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"1234": "foo",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathConversionError)
}

func TestJSONPathNotEqual(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()

	exp := gdtjson.Expect{
		Paths: map[string]string{
			"$[0].pages": "127",
		},
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK(ctx))
	require.Empty(a.Failures())

	exp = gdtjson.Expect{
		Paths: map[string]string{
			"$[0].pages": "42",
		},
	}

	a = gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathNotEqual)
}

func TestJSONPathFormatNotFound(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()

	exp := gdtjson.Expect{
		PathFormats: map[string]string{
			"$.noexist": "invalidformat",
		},
	}

	a := gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONPathNotFound)
}

func TestJSONPathFormatNotEqual(t *testing.T) {
	require := require.New(t)

	ctx := context.TODO()
	c := content()

	exp := gdtjson.Expect{
		PathFormats: map[string]string{
			"$[0].id": "uuid4",
		},
	}

	a := gdtjson.New(&exp, c)
	require.True(a.OK(ctx))
	require.Empty(a.Failures())

	exp = gdtjson.Expect{
		PathFormats: map[string]string{
			"$[0].pages": "uuid4",
		},
	}

	a = gdtjson.New(&exp, c)
	require.False(a.OK(ctx))
	failures := a.Failures()
	require.Len(failures, 1)
	require.ErrorIs(failures[0], gdtjson.ErrJSONFormatNotEqual)
}
