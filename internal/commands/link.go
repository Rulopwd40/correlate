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
	Long:    "Add a project reference",
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
