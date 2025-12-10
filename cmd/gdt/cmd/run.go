package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdt-dev/core/api"
	gdtcontext "github.com/gdt-dev/core/context"
	"github.com/gdt-dev/core/run"
	"github.com/gdt-dev/core/scenario"
	"github.com/gdt-dev/core/suite"
	"github.com/gdt-dev/core/testunit"
	"github.com/spf13/cobra"

	"github.com/gdt-dev/gdt/cmd/gdt/pkg/cli"
)

const (
	debugPrefix = "[gdt]"
	runUsage    = `run <subject> [<subject> ...]`
	runDescLong = `Check test scenarios or test suites for parse errors.

The command will run gdt test scenarios or test suites pointed to by <subject>.

<subject> should be a path to a YAML file or a directory containing YAML files.

Returns 0 if all subject test scenarios complete without failure, 1 otherwise.
`
)

var (
	runOutputFormatHuman   = "human"
	runOutputFormatJSON    = "json"
	runOutputFormatXUnit   = "xunit"
	defaultRunOutputFormat = runOutputFormatHuman
	supportedOutputFormats = []string{
		runOutputFormatHuman,
		runOutputFormatJSON,
		runOutputFormatXUnit,
	}
	optRunOutputFormat   = defaultRunOutputFormat
	usageRunOutputFormat = `output format ("human","json","xunit")`
)

var RunCmd = &cobra.Command{
	Use:     runUsage,
	Short:   "run test scenario/suites.",
	Long:    runDescLong,
	Aliases: []string{"exec"},
	RunE:    doRun,
}

func init() {
	RunCmd.Flags().BoolVarP(
		&optQuiet,
		"quiet",
		"q",
		false,
		optQuietUsage,
	)
	RunCmd.Flags().StringVarP(
		&optRunOutputFormat,
		"output-format", "o",
		defaultRunOutputFormat, usageRunOutputFormat,
	)
}

func doRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("supply <subject> containing filepath to YAML file or directory.")
	}
	if cli.CommonOptions.Debug {
		cli.CommonOptions.Verbose = true
	}
	ctx := gdtcontext.New(gdtcontext.WithDebugPrefix(debugPrefix))
	runs := []*run.Run{}
	for _, path := range args {
		fi, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%q not found.", path)
			}
			return err
		}
		run := run.New()
		if fi.IsDir() {
			cli.Df("loading suite from directory %q ...", path)
			su, err := suite.FromDir(path)
			if err != nil {
				return err
			}
			err = su.Run(ctx, run)
			if err != nil {
				// Run() only returns RuntimeErrors. The `run` object will
				// contain assertion failures, which are not considered
				// RuntimeErrors.
				return err
			}
		} else {
			cli.Df("loading scenario from file %q ...", path)
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			sc, err := scenario.FromReader(f, scenario.WithPath(path))
			if err != nil {
				return err
			}
			err = sc.Run(ctx, run)
			if err != nil {
				// Run() only returns RuntimeErrors. The `run` object will
				// contain assertion failures, which are not considered
				// RuntimeErrors.
				return err
			}
		}
		// We do this here so that we get more immediate output results
		// when using human output format
		if optRunOutputFormat == runOutputFormatHuman {
			printRun(run)
		}
		runs = append(runs, run)
	}
	if optRunOutputFormat != runOutputFormatHuman {
		return printResults(runs)
	}
	return nil
}

// printRun outputs the human-readable results of a test scenario run.
func printRun(
	run *run.Run,
) {
	if !optQuiet {
		paths := run.ScenarioPaths()
		for _, path := range paths {
			shortPath := filepath.Base(path)
			if cli.CommonOptions.Verbose {
				fmt.Printf("=== RUN: %s\n", shortPath)
			}
			var scenElapsed time.Duration

			results := run.ScenarioResults(path)
			scenOK := true

			for _, res := range results {
				scenElapsed += res.Elapsed()
				scenOK = scenOK && res.OK()
				printTestUnitResult(res)
			}

			if scenOK {
				if cli.CommonOptions.Verbose {
					fmt.Printf("PASS (%s)\n", scenElapsed)
				} else {
					fmt.Printf("ok\t%s\t%s\n", path, scenElapsed)
				}
			} else {
				fmt.Printf("FAIL\t%s\t%s\n", path, scenElapsed)
			}
		}
	}
	if !run.OK() {
		if cli.CommonOptions.Verbose {
			fmt.Println("FAIL")
		}
		os.Exit(1)
	} else {
		if cli.CommonOptions.Verbose {
			fmt.Println("PASS")
		}
	}
}

func printTestUnitResult(r testunit.Result) {
	if r.Skipped() {
		if cli.CommonOptions.Verbose {
			fmt.Printf("--- SKIP: %s (%s)\n", r.Name(), r.Elapsed())
		}
	} else if r.OK() {
		if cli.CommonOptions.Verbose {
			fmt.Printf("--- PASS: %s (%s)\n", r.Name(), r.Elapsed())
		}
	} else {
		for _, fail := range r.Failures() {
			indentFail := indent(fail.Error(), 1)
			if !optQuiet {
				fmt.Printf(
					"--- FAIL: %s (%s)\n%s\n",
					r.Name(), r.Elapsed(), indentFail,
				)
			}
		}
	}

	if cli.CommonOptions.Debug || !r.OK() {
		detail := r.Detail()
		if len(detail) > 0 {
			cli.HorizontalSectionHeader("detail")
			fmt.Printf("%s", r.Detail())
			cli.HorizontalBar()
		}
	}
}

// printResults outputs the test results of all test scenarios for JSON or
// XUnit output.
func printResults(
	runs []*run.Run,
) error {
	suites := []api.XUnitTestSuite{}
	for _, r := range runs {
		suites = append(suites, r.XUnit()...)
	}
	res := api.XUnitResults{
		TestSuites: suites,
	}
	var err error
	var out string
	var b []byte
	if optRunOutputFormat == runOutputFormatJSON {
		b, err = json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		out = string(b)
	} else {
		b, err = xml.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		out = xml.Header + string(b)
	}
	fmt.Println(out)
	return nil
}

func indent(subject string, level int) string {
	indentStr := strings.Repeat(" ", level*4)
	b := strings.Builder{}
	lines := strings.Split(subject, "\n")
	for _, line := range lines {
		b.WriteString(fmt.Sprintf("%s%s", indentStr, line))
	}
	return b.String()
}
