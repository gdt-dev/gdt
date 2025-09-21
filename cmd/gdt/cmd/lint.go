package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdt-dev/core/parse"
	"github.com/gdt-dev/core/scenario"
	_ "github.com/gdt-dev/kube"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/gdt-dev/gdt/cmd/gdt/pkg/cli"
)

const (
	lintUsage    = `lint <subject> [<subject> ...]`
	lintDescLong = `Check test scenarios or test suites for parse errors.

The command will load gdt test scenarios or test suites pointed to by <subject>
and report any parse errors.

<subject> should be a path to a YAML file or a directory containing YAML files.

Returns 0 on successful parsing of all subjects, 1 if any parse errors are
detected.
`
	optQuietUsage = `quiet output. only return an exit code 0 for success, 1 for failure.`
)

var (
	optQuiet = false
)

var LintCmd = &cobra.Command{
	Use:     lintUsage,
	Short:   "check test scenario/suites for parse errors.",
	Long:    lintDescLong,
	Aliases: []string{"check", "parse"},
	RunE:    doLint,
}

func init() {
	LintCmd.Flags().BoolVarP(
		&optQuiet,
		"quiet",
		"q",
		false,
		optQuietUsage,
	)
}

func doLint(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("supply <subject> containing filepath to YAML file or directory.")
	}
	results := []lintResult{}
	for _, path := range args {
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%q not found.", path)
			}
			return err
		}
		if fi.IsDir() {
			cli.Df("checking directory %q ...", path)
			dirResults, err := lintDir(path)
			if err != nil {
				return err
			}
			results = append(results, dirResults...)
		} else {
			cli.Df("checking file %q ...", path)
			res := lintResult{path: path}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			// Need to chdir here so that test scenario may reference files in
			// relative directories
			if err := os.Chdir(filepath.Dir(path)); err != nil {
				return err
			}

			sc, err := scenario.FromReader(f, scenario.WithPath(path))
			if err != nil {
				if ep, ok := err.(*parse.Error); ok {
					res.err = ep
				} else {
					return err
				}
			}
			res.scenario = sc
			results = append(results, res)
		}
	}

	if !optQuiet {
		for _, res := range results {
			res.print()
		}
	}
	if lo.SomeBy(results, func(r lintResult) bool {
		return r.err != nil
	}) {
		os.Exit(1)
	}
	return nil
}

var (
	validFileExts = []string{".yaml", ".yml"}
)

// lintDir reads the supplied directory path and parses any YAML files
// contained within it.
func lintDir(
	dirPath string,
) ([]lintResult, error) {
	if _, err := os.Stat(dirPath); err != nil {
		return nil, err
	}
	results := []lintResult{}
	if err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, _ error) error {
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

			res := lintResult{path: path}
			sc, err := scenario.FromReader(f, scenario.WithPath(path))
			if err != nil {
				if ep, ok := err.(*parse.Error); ok {
					res.err = ep
				} else {
					return err
				}
			}
			res.scenario = sc
			results = append(results, res)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return results, nil
}

type lintResult struct {
	path     string
	scenario *scenario.Scenario
	err      *parse.Error
}

func (r lintResult) shortPath() string {
	p := r.path
	for _, ext := range validFileExts {
		p = strings.TrimSuffix(p, ext)
	}
	parts := strings.Split(p, string(os.PathSeparator))
	if len(parts) == 1 {
		return parts[0]
	}
	return fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
}

func (r lintResult) print() {
	if r.err == nil {
		fmt.Printf("%s: ok.\n", r.shortPath())
	} else {
		fmt.Printf("%s: fail\n", r.shortPath())
		cli.HorizontalBar()
		fmt.Printf("%s", r.err)
		cli.HorizontalBar()
	}
}
