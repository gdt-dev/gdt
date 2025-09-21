package main

import (
	"os"

	"github.com/spf13/cobra"

	_ "github.com/gdt-dev/core/plugin/exec"
	"github.com/gdt-dev/gdt/cmd/gdt/cmd"
	"github.com/gdt-dev/gdt/cmd/gdt/pkg/cli"
	_ "github.com/gdt-dev/kube"
)

var rootCmd = &cobra.Command{
	Use:   "gdt",
	Short: "gdt - declarative functional testing",
	Long: `           __ __   
 .-----.--|  |  |_ 
 |  _  |  _  |   _|
 |___  |_____|____|
 |_____|

Functional testing that makes sense.

https://github.com/gdt-dev/gdt
`,
}

func init() {
	rootCmd.PersistentFlags().AddFlagSet(cli.CommonOptionsFlagSet)

	rootCmd.AddCommand(cmd.LintCmd)
	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.SilenceUsage = true
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
