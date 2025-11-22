package commands

import (
	"log"

	"github.com/Rulopwd40/correlate/internal/core"
	"github.com/spf13/cobra"
)

var version string

var ReplaceCmd = &cobra.Command{
	Use:   "replace [identifier]",
	Short: "Replace version in project dependency",
	Args:  cobra.ExactArgs(1),
	Long:  "Replace version in project dependency. If --version is not specified, it will search for it in the target manifest.",
	Run: func(cmd *cobra.Command, args []string) {
		runReplace(cmd, args)
	},
}

func init() {
	ReplaceCmd.Flags().StringVarP(&version, "version", "v", "", "Version to apply (optional)")
}

func runReplace(cmd *cobra.Command, args []string) {
	identifier := args[0]
	orch, err := core.Get[*core.Orchestrator]()
	if err != nil {
		log.Println("Error getting orchestrator:", err)
		return
	}
	err = orch.Replace(identifier, version)
	if err != nil {
		log.Println("Error replacing version in project:", err)
		return
	}

	log.Println("Successfully replaced version in project.")
}
