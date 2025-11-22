package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "correlate",
	Short: "Correlate automation CLI",
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(LinkCmd)
	RootCmd.AddCommand(ReplaceCmd)
	RootCmd.AddCommand(UpdateCmd)
}
