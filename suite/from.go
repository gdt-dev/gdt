// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"os"
	"path/filepath"

	"github.com/gdt-dev/gdt/scenario"
	"github.com/samber/lo"
)

var (
	validFileExts = []string{".yaml", ".yml"}
)

// FromDir reads the supplied directory path and returns a Suite representing
// the suite of test scenarios in that directory.
func FromDir(
	dirPath string,
	mods ...SuiteModifier,
) (*Suite, error) {
	if _, err := os.Stat(dirPath); err != nil {
		return nil, err
	}
	// List YAML files in the directory and parse each into a testable unit
	mods = append(mods, WithPath(dirPath))
	s := New(mods...)

	if err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			suffix := filepath.Ext(path)
			if !lo.Contains(validFileExts, suffix) {
				return nil
			}
			f, err := os.Open(path)

			if err != nil {
				return err
			}
			defer f.Close()

			tc, err := scenario.FromReader(f, scenario.WithPath(path))
			if err != nil {
				return err
			}
			s.Append(tc)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return s, nil
}
