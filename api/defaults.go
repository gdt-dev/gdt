// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

// Defaults are a collection of default configuration values
type Defaults map[string]interface{}

func (d *Defaults) For(key string) interface{} {
	if d == nil {
		return nil
	}
	return (*d)[key]
}
