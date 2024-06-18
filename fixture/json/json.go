// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	gdttypes "github.com/gdt-dev/gdt/types"
)

type jsonFixture struct {
	data interface{}
}

func (f *jsonFixture) Start(_ context.Context) error { return nil }

func (f *jsonFixture) Stop(_ context.Context) {}

// HasState returns true if the supplied JSONPath expression results in a found
// value in the fixture's data
func (f *jsonFixture) HasState(path string) bool {
	if f.data == nil {
		return false
	}
	got, err := jsonpath.Get(path, f.data)
	if err != nil {
		return false
	}
	if got == nil {
		return false
	}
	return true
}

// GetState returns the value at supplied JSONPath expression or nil if the
// JSONPath expression does not result in any matched field
func (f *jsonFixture) State(path string) interface{} {
	if f.data == nil {
		return nil
	}
	got, err := jsonpath.Get(path, f.data)
	if err != nil {
		return nil
	}
	switch got.(type) {
	case string:
		return got.(string)
	case float64:
		return strconv.FormatFloat(got.(float64), 'f', 0, 64)
	default:
		return nil
	}
}

// New takes a string, some bytes or an io.Reader and returns a new
// gdttypes.Fixture that can have its state queried via JSONPath
func New(data interface{}) (gdttypes.Fixture, error) {
	var err error
	var b []byte
	switch data.(type) {
	case io.Reader:
		r := data.(io.Reader)
		b, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	case []byte:
		b = data.([]byte)
	case string:
		b = []byte(data.(string))
	}
	f := jsonFixture{
		data: interface{}(nil),
	}
	if err = json.Unmarshal(b, &f.data); err != nil {
		return nil, err
	}
	return &f, nil
}
