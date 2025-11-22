package commands

import (
	"log"

	"github.com/Rulopwd40/correlate/internal/core"
	"github.com/spf13/cobra"
)

var LinkCmd = &cobra.Command{
	Use:     "link [identifier] [fullPath]",
	Aliases: []string{"l"},
	Short:   "Add a project reference",
	Args:    cobra.ExactArgs(2),
	Long: `Link a dependency to your correlate project.

This command scans the specified project path for all occurrences of the
identifier in manifest files (e.g., pom.xml) and adds them to the references
file for tracking and automated updates.

Examples:
  correlate link my-library /path/to/project
  correlate link g10-deliverymanagementsystem C:\\Projects\\consumer-app
  correlate l my-library ../dependent-project`,
	Run: func(cmd *cobra.Command, args []string) {
		runLink(cmd, args)
	},
}

func runLink(cmd *cobra.Command, args []string) {
	log.Println("Linking project...")
	orch, err := core.Get[*core.Orchestrator]()
	if err != nil {
		log.Println("Error getting orchestrator:", err)
	}
	identifier := args[0]
	fullPath := args[1]

	err = orch.Link(identifier, fullPath)
	if err != nil {
		return
	}
	log.Println("Project successfully linked!")
}
