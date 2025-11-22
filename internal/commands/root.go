package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "correlate",
	Short: "Correlate - Dependency management automation tool",
	Long: `Correlate is a CLI tool for managing and automating updates across
multiple projects with interdependencies.

It helps you:
  - Track dependencies between projects
  - Automatically detect and update versions
  - Execute build and update pipelines
  - Maintain consistency across multiple repositories

Use "correlate [command] --help" for more information about a command.`,
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(LinkCmd)
	RootCmd.AddCommand(ReplaceCmd)
	RootCmd.AddCommand(UpdateCmd)
	RootCmd.AddCommand(VersionCmd)
}
