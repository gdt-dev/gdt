// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package json

import (
	"fmt"

	"github.com/google/uuid"
	gjs "github.com/xeipuuv/gojsonschema"
)

var (
	validators = map[string]gjs.FormatChecker{
		"date":                  gjs.DateFormatChecker{},
		"time":                  gjs.TimeFormatChecker{},
		"date-time":             gjs.DateTimeFormatChecker{},
		"hostname":              gjs.HostnameFormatChecker{},
		"email":                 gjs.EmailFormatChecker{},
		"idn-email":             gjs.EmailFormatChecker{},
		"ipv4":                  gjs.IPV4FormatChecker{},
		"ipv6":                  gjs.IPV6FormatChecker{},
		"uri":                   gjs.URIFormatChecker{},
		"uri-reference":         gjs.URIReferenceFormatChecker{},
		"iri":                   gjs.URIFormatChecker{},
		"iri-reference":         gjs.URIReferenceFormatChecker{},
		"uri-template":          gjs.URITemplateFormatChecker{},
		"uuid":                  gjs.UUIDFormatChecker{},
		"regex":                 gjs.RegexFormatChecker{},
		"json-pointer":          gjs.JSONPointerFormatChecker{},
		"relative-json-pointer": gjs.RelativeJSONPointerFormatChecker{},
		"uuid4":                 uuid4FormatChecker{},
	}
)

// isFormatted takes a format string and a string value, determines the
// validator function for that type of format string and returns whether the
// value string is formatted correctly.
func isFormatted(format string, input interface{}) (bool, error) {
	c, ok := validators[format]
	if !ok {
		return false, fmt.Errorf("unknown format %s", format)
	}
	return c.IsFormat(input), nil
}

type uuid4FormatChecker struct{}

func (c uuid4FormatChecker) IsFormat(input interface{}) bool {
	s, ok := input.(string)
	if !ok {
		return false
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == 4
}
