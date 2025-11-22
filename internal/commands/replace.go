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
	Long: `Replace the version of a dependency in the project manifest.

If --version is not specified, correlate will automatically detect the version
from the linked project's manifest. This ensures version consistency across
all dependent projects.

Examples:
  correlate replace my-library --version 1.2.3
  correlate replace my-library -v 2.0.0
  correlate replace my-library  # Auto-detect version`,
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
