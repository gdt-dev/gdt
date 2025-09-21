package cli

import "github.com/spf13/pflag"

type commonOpts struct {
	Debug   bool
	Verbose bool
}

var (
	CommonOptions commonOpts
)

var (
	CommonOptionsFlagSet = &pflag.FlagSet{}
)

func init() {
	CommonOptionsFlagSet.BoolVarP(
		&CommonOptions.Debug,
		"debug", "d", false,
		"Enable or disable debug mode",
	)
	CommonOptionsFlagSet.BoolVarP(
		&CommonOptions.Verbose,
		"verbose", "v",
		false,
		"Enable or disable verbose output.",
	)
}
