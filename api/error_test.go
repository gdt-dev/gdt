// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api_test

import (
	"testing"

	"github.com/gdt-dev/gdt/api"
	"github.com/stretchr/testify/assert"
)

func TestUnknownSourceType(t *testing.T) {
	assert := assert.New(t)

	err := api.UnknownSourceType(1)
	assert.ErrorContains(err, "int")

	source := []string{"foo", "bar"}
	err = api.UnknownSourceType(source)
	assert.ErrorContains(err, "[]string")
}
